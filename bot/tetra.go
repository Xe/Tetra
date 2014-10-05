/*
Package tetra implements the core for a TS6 pseudoserver. It also has lua and
moonscript loading support to add functionality at runtime.
*/
package tetra

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"log/syslog"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/atheme"
	"github.com/coreos/go-etcd/etcd"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/influxdb"
	"github.com/robfig/cron"
)

// Struct Tetra contains all fields for Tetra.
type Tetra struct {
	Conn     *Connection
	Info     *Server
	Clients  *Clients
	Channels map[string]*Channel
	Bursted  bool
	Handlers map[string]map[string]*Handler
	Services map[string]*Client
	Servers  map[string]*Server
	Scripts  map[string]*Script
	Hooks    map[string][]*Hook
	nextuid  int
	Config   *Config
	Log      *log.Logger
	Uplink   *Server
	tasks    chan string
	wg       *sync.WaitGroup
	Etcd     *etcd.Client
	Atheme   *atheme.Atheme
	Cron     *cron.Cron
}

// NewTetra returns a new instance of Tetra based on a config file located at cpath.
// This also kicks off the worker goroutines and statistics collection, as well
// as seeding basic protocol verb handlers.
func NewTetra(cpath string) (tetra *Tetra) {
	config, err := NewConfig(cpath)
	if err != nil {
		fmt.Printf("No config file %s\n", cpath)
		panic(err)
	}

	if config.General.Workers == 0 {
		config.General.Workers = 4
	}

	tetra = &Tetra{
		Conn: &Connection{
			Log:    log.New(os.Stdout, "CONN ", log.LstdFlags),
			Buffer: make(chan string, 100),
			open:   true,
		},
		Info: &Server{},
		Clients: &Clients{
			ByNick: make(map[string]*Client),
			ByUID:  make(map[string]*Client),
			Tetra:  tetra,
			Gauge:  metrics.NewGauge(),
		},
		Channels: make(map[string]*Channel),
		Handlers: make(map[string]map[string]*Handler),
		Services: make(map[string]*Client),
		Servers:  make(map[string]*Server),
		Scripts:  make(map[string]*Script),
		Hooks:    make(map[string][]*Hook),
		Bursted:  false,
		nextuid:  100000,
		Config:   config,
		Log:      log.New(os.Stdout, "BOT ", log.LstdFlags),
		Uplink: &Server{
			Counter: metrics.NewGauge(),
		},
		tasks: make(chan string, 100),
		wg:    &sync.WaitGroup{},
		Cron:  cron.New(),
	}

	tetra.Info = &Server{
		Sid:     tetra.Config.Server.Sid,
		Name:    tetra.Config.Server.Name,
		Gecos:   tetra.Config.Server.Gecos,
		Links:   []*Server{tetra.Uplink},
		Counter: nil,
	}

	tetra.Conn.Debug = tetra.Config.General.Debug

	tetra.Etcd = etcd.NewClient(tetra.Config.Etcd.Machines)
	tetra.Etcd.CreateDir("/tetra", 0)
	tetra.Etcd.CreateDir("/tetra/channels", 0)
	tetra.Etcd.CreateDir("/tetra/clients", 0)
	tetra.Etcd.CreateDir("/tetra/scripts", 0)

	tetra.Atheme, err = atheme.NewAtheme(tetra.Config.Atheme.URL)
	if err != nil {
		tetra.Log.Fatal(err)
	}

	if tetra.Atheme == nil {
		tetra.Log.Fatal("tetra.Atheme is nil.")
	}

	err = tetra.Atheme.Login(tetra.Config.Atheme.Username, tetra.Config.Atheme.Password)

	tetra.Cron.AddFunc("0 30 * * * *", func() {
		debug("Keeping us logged into Atheme...")
		tetra.Atheme.MemoServ.List()
	})

	tetra.Cron.AddFunc("@every 5m", func() {
		tetra.RunHook("CRON-HEARTBEAT")
	})

	tetra.NewHook("CRON-HEARTBEAT", tetra.GetChannelStats)
	tetra.NewHook("CRON-HEARTBEAT", tetra.GetChannelStats)

	tetra.Cron.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			_ = sig

			fmt.Println(" <-- Control-C pressed!")

			tetra.Quit()
		}
	}()

	tetra.seedHandlers()

	metrics.Register(tetra.Config.Server.Name+"_clients", tetra.Info.Counter)

	tetra.startWorkers(config.General.Workers)

	go tetra.Conn.sendLinesWait()

	return
}

// NextUID returns a new TS6 UID.
func (tetra *Tetra) NextUID() string {
	tetra.nextuid++
	return tetra.Info.Sid + strconv.Itoa(tetra.nextuid)
}

// Connect connects to the uplink server.
func (tetra *Tetra) Connect(host, port string) (err error) {
	if tetra.Config.Uplink.Ssl {
		config := &tls.Config{InsecureSkipVerify: true}
		tetra.Conn.Conn, err = tls.Dial("tcp", host+":"+port, config)
	} else {
		tetra.Conn.Conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			tetra.Log.Fatal(err)
		}
	}

	tetra.Conn.Reader = bufio.NewReader(tetra.Conn.Conn)
	tetra.Conn.Tp = textproto.NewReader(tetra.Conn.Reader)

	return
}

// Auth authenticates over TS6.
func (tetra *Tetra) Auth() {
	tetra.Conn.SendLine("PASS " + tetra.Config.Uplink.Password + " TS 6 :" + tetra.Config.Server.Sid)
	tetra.Conn.SendLine("CAPAB :QS EX IE KLN UNKLN ENCAP SERVICES EUID EOPMO")
	tetra.Conn.SendLine("SERVER " + tetra.Config.Server.Name + " 1 :" + tetra.Config.Server.Gecos)
}

// Burst sends our local information after recieving the server's burst.
func (tetra *Tetra) Burst() {
	for _, script := range tetra.Config.Autoload {
		tetra.LoadScript(script)
	}

	for _, client := range tetra.Services {
		tetra.Conn.SendLine(client.Euid())
		if client.Certfp != "" {
			tetra.Conn.SendLine(":%s ENCAP * CERTFP :%s", client.Uid, client.Certfp)
		}

		client.Join(tetra.Config.General.SnoopChan)
	}

	for _, channel := range tetra.Channels {
		for uid, _ := range channel.Clients {
			if !strings.HasPrefix(uid, tetra.Info.Sid) {
				continue
			}
			str := fmt.Sprintf(":%s SJOIN %d %s + :%s", tetra.Info.Sid, channel.Ts, channel.Name, uid)
			tetra.Conn.SendLine(str)
		}
	}

	metrics.Register("clientcount", tetra.Clients.Gauge)

	if tetra.Config.Stats.Host != "NOCOLLECTION" {
		w, _ := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
		go metrics.Syslog(metrics.DefaultRegistry, 60e9, w)

		go tetra.GetNetworkStats()
		go tetra.GetChannelStats()

		go influxdb.Influxdb(metrics.DefaultRegistry, 5*time.Minute, &influxdb.Config{
			Host:     tetra.Config.Stats.Host,
			Database: tetra.Config.Stats.Database,
			Username: tetra.Config.Stats.Username,
			Password: tetra.Config.Stats.Password,
		})
	}

	tetra.Bursted = true
}

// StickConfig creates Clients based off of the config file and handles module
// autoloads.
func (tetra *Tetra) StickConfig() {
	for _, sclient := range tetra.Config.Services {
		client := tetra.AddService(sclient.Name, sclient.Nick, sclient.User, sclient.Host, sclient.Gecos, sclient.Certfp)

		filepath.Walk("modules/"+client.Kind+"/core/", func(path string, info os.FileInfo, err error) error {
			modname := strings.Split(path, ".")[0]
			mods := strings.Split(modname, "/")
			modname = mods[len(mods)-1]

			if len(modname) == 0 {
				return nil
			}

			tetra.LoadScript(client.Kind + "/core/" + modname)

			return nil
		})

		client.NewCommand("HELP", func(source *Client, target Targeter, message []string) (ret string) {
			if len(message) == 0 {
				if helpHas(client.Kind, "_index") {
					client.showHelp(source, client.Kind, "_index")

					return "End of help file"
				} else {
					return "No help available."
				}
			}

			basecommand := strings.ToUpper(message[0])

			command := strings.ToLower(strings.Join(message, " "))

			if _, present := client.Commands[basecommand]; !present {
				if helpHas(client.Kind, command) {
					client.showHelp(source, client.Kind, command)

					return "End of help file"
				}
			}

			if helpHas(client.Kind, command) {
				if client.Commands[basecommand].NeedsOper && !source.IsOper() {
					return "Permission denied."
				} else {
					client.showHelp(source, client.Kind, command)

					return "End of help file"
				}
			} else {
				if _, present := client.Commands[basecommand]; !present {
					return "No such command " + basecommand
				}
			}

			return "Help for " + strings.ToUpper(command) + " not found."
		})
	}
}

// Quit kills Tetra gracefully.
func (tetra *Tetra) Quit() {
	for _, service := range tetra.Services {
		tetra.DelService(service.Kind)
	}

	tetra.Conn.SendLine("SQUIT %s :Goodbye!", tetra.Info.Sid)
}

// ProcessLine processes a line as if it came from the server.
func (tetra *Tetra) ProcessLine(line string) {
	rawline := r1459.NewRawLine(line)

	debugf("<<< %s", line)

	// This should just be hard-coded here.
	if rawline.Verb == "PING" {
		if rawline.Source == "" {
			if !tetra.Bursted {
				tetra.Burst()
				tetra.Log.Printf("Bursted!")
			}
			tetra.Conn.SendLine("PONG :" + rawline.Args[0])
		} else {
			tetra.Conn.SendLine(":%s PONG %s :%s", tetra.Info.Sid, tetra.Info.Name, rawline.Source)
		}
	}

	if _, present := tetra.Handlers[rawline.Verb]; present {
		for _, handler := range tetra.Handlers[rawline.Verb] {
			defer func() {
				if r := recover(); r != nil {
					str := fmt.Sprintf("Recovered in handler %s (%s): %#v",
						handler.Verb, handler.Uuid, r)
					tetra.Log.Print(str)
					tetra.Log.Printf("%#v", r)
					tetra.Services["tetra"].ServicesLog(str)
				}
			}()
			handler.Impl(rawline)
		}
	}
}

// Main is the main loop.
func (t *Tetra) Main() {
	for {
		line, err := t.Conn.GetLine()
		if err != nil {
			break
		}

		debug("Got line")

		if t.Bursted {
			t.tasks <- line
		} else {
			debug("begin process line")
			t.ProcessLine(line)
			debug("End process line")
		}
	}

	t.wg.Wait()
}
