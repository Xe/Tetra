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
	"github.com/Xe/Tetra/bot/config"
	"github.com/coreos/go-etcd/etcd"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/influxdb"
	"github.com/robfig/cron"
)

// Struct Tetra contains all fields for
var (
	Conn         *Connection
	Info         *Server
	Clients      *ClientSet
	Channels     map[string]*Channel
	Bursted      bool
	Handlers     map[string]map[string]*Handler
	Services     map[string]*Client
	Servers      map[string]*Server
	Scripts      map[string]*Script
	Hooks        map[string][]*Hook
	nextuid      int
	ActiveConfig *config.Config
	Log          *log.Logger
	Uplink       *Server
	tasks        chan string
	wg           *sync.WaitGroup
	Etcd         *etcd.Client
	Atheme       *atheme.Atheme
	Cron         *cron.Cron
	lock         *sync.Mutex
)

// NewTetra returns a new instance of Tetra based on a config file located at cpath.
// This also kicks off the worker goroutines and statistics collection, as well
// as seeding basic protocol verb handlers.
func NewTetra(cpath string) {
	config, err := config.NewConfig(cpath)
	if err != nil {
		fmt.Printf("No config file %s\n", cpath)
		panic(err)
	}

	if config.General.Workers == 0 {
		config.General.Workers = 4
	}

	Conn = &Connection{
		Log:    log.New(os.Stdout, "CONN ", log.LstdFlags),
		Buffer: make(chan string, 100),
		open:   true,
	}
	Info = &Server{}
	Clients = &ClientSet{
		ByNick: make(map[string]*Client),
		ByUID:  make(map[string]*Client),
		Gauge:  metrics.NewGauge(),
	}
	Channels = make(map[string]*Channel)
	Handlers = make(map[string]map[string]*Handler)
	Services = make(map[string]*Client)
	Servers = make(map[string]*Server)
	Scripts = make(map[string]*Script)
	Hooks = make(map[string][]*Hook)
	Bursted = false
	nextuid = 100000
	ActiveConfig = config
	Log = log.New(os.Stdout, "BOT ", log.LstdFlags)
	Uplink = &Server{
		Counter: metrics.NewGauge(),
	}
	tasks = make(chan string, 100)
	wg = &sync.WaitGroup{}
	Cron = cron.New()

	Info = &Server{
		Sid:     ActiveConfig.Server.Sid,
		Name:    ActiveConfig.Server.Name,
		Gecos:   ActiveConfig.Server.Gecos,
		Links:   []*Server{Uplink},
		Counter: nil,
	}

	Conn.Debug = ActiveConfig.General.Debug

	Etcd = etcd.NewClient(ActiveConfig.Etcd.Machines)
	Etcd.CreateDir("/tetra", 0)
	Etcd.CreateDir("/tetra/channels", 0)
	Etcd.CreateDir("/tetra/clients", 0)
	Etcd.CreateDir("/tetra/scripts", 0)

	Atheme, err = atheme.NewAtheme(ActiveConfig.Atheme.URL)
	if err != nil {
		Log.Fatal(err)
	}

	if Atheme == nil {
		Log.Fatal("Atheme is nil.")
	}

	err = Atheme.Login(ActiveConfig.Atheme.Username, ActiveConfig.Atheme.Password)
	if err != nil {
		Log.Fatalf("Atheme error: %s", err.Error())
	}

	Cron.AddFunc("0 30 * * * *", func() {
		debug("Keeping us logged into Atheme...")
		_, err := Atheme.MemoServ.List()
		if err != nil {
			err = Atheme.Login(ActiveConfig.Atheme.Username, ActiveConfig.Atheme.Password)
			if err != nil {
				Log.Fatalf("Atheme error: %s", err.Error())
			}
		}
	})

	Cron.AddFunc("@every 5m", func() {
		RunHook("CRON-HEARTBEAT")
	})

	Cron.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			_ = sig

			fmt.Println(" <-- Control-C pressed!")

			Quit()
		}
	}()

	seedHandlers()

	startWorkers(config.General.Workers)

	lock = &sync.Mutex{}

	go Conn.sendLinesWait()

	return
}

// NextUID returns a new TS6 UID.
func NextUID() string {
	nextuid++
	return Info.Sid + strconv.Itoa(nextuid)
}

// Connect connects to the uplink server.
func Connect(host, port string) (err error) {
	if ActiveConfig.Uplink.Ssl {
		config := &tls.Config{InsecureSkipVerify: true}
		Conn.Conn, err = tls.Dial("tcp", host+":"+port, config)
	} else {
		Conn.Conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			Log.Fatal(err)
		}
	}

	Conn.Reader = bufio.NewReader(Conn.Conn)
	Conn.Tp = textproto.NewReader(Conn.Reader)

	return
}

// Auth authenticates over TS6.
func Auth() {
	Conn.SendLine("PASS " + ActiveConfig.Uplink.Password + " TS 6 :" + ActiveConfig.Server.Sid)
	Conn.SendLine("CAPAB :QS EX IE KLN UNKLN ENCAP SERVICES EUID RSFNC SAVE MLOCK CHW TB CLUSTER BAN")
	Conn.SendLine("SERVER " + ActiveConfig.Server.Name + " 1 :" + ActiveConfig.Server.Gecos)
}

// Burst sends our local information after recieving the server's burst.
func Burst() {
	for _, script := range ActiveConfig.Autoload {
		LoadScript(script)
	}

	for _, client := range Services {
		Conn.SendLine(client.Euid())
		if client.Certfp != "" {
			Conn.SendLine(":%s ENCAP * CERTFP :%s", client.Uid, client.Certfp)
		}

		client.Join(ActiveConfig.General.SnoopChan)
	}

	for _, channel := range Channels {
		for uid := range channel.Clients {
			if !strings.HasPrefix(uid, Info.Sid) {
				continue
			}
			str := fmt.Sprintf(":%s SJOIN %d %s + :%s", Info.Sid, channel.Ts, channel.Name, uid)
			Conn.SendLine(str)
		}
	}

	metrics.Register("clientcount", Clients.Gauge)

	if ActiveConfig.Stats.Host != "NOCOLLECTION" {
		go GetNetworkStats()
		go GetChannelStats()

		go influxdb.Influxdb(metrics.DefaultRegistry, 5*time.Minute, &influxdb.Config{
			Host:     ActiveConfig.Stats.Host,
			Database: ActiveConfig.Stats.Database,
			Username: ActiveConfig.Stats.Username,
			Password: ActiveConfig.Stats.Password,
		})

		metrics.Register(ActiveConfig.Server.Name+"_clients", Info.Counter)
	}

	Bursted = true

	RunHook("BURSTED")
}

// StickConfig creates Clients based off of the config file and handles module
// autoloads.
func StickConfig() {
	for _, sclient := range ActiveConfig.Services {
		client := AddService(sclient.Name, sclient.Nick, sclient.User, sclient.Host, sclient.Gecos, sclient.Certfp)

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

	time.Sleep(500 * time.Millisecond)

	for _, client := range Services {
		filepath.Walk("modules/"+client.Kind+"/core/", func(path string, info os.FileInfo, err error) error {
			modname := strings.Split(path, ".")[0]
			mods := strings.Split(modname, "/")
			modname = mods[len(mods)-1]

			if len(modname) == 0 {
				return nil
			}

			LoadScript(client.Kind + "/core/" + modname)

			return nil
		})
	}
}

// Quit kills Tetra gracefully.
func Quit() {
	for _, service := range Services {
		DelService(service.Kind)
	}

	Conn.SendLine("SQUIT %s :Goodbye!", Info.Sid)
}

// ProcessLine processes a line as if it came from the server.
func ProcessLine(line string) {
	rawline := r1459.NewRawLine(line)

	debugf("<<< %s", line)

	// This should just be hard-coded here.
	if rawline.Verb == "PING" {
		if rawline.Source == "" {
			if !Bursted {
				Burst()
				Log.Printf("Bursted!")
			}
			Conn.SendLine("PONG :" + rawline.Args[0])
		} else {
			Conn.SendLine(":%s PONG %s :%s", Info.Sid, Info.Name, rawline.Source)
		}
	}

	if _, present := Handlers[rawline.Verb]; present {
		for _, handler := range Handlers[rawline.Verb] {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						str := fmt.Sprintf("Recovered in handler for %s %#v (%s), sleeping and retrying",
							handler.Verb, r, line)
						Log.Print(str)
					} else {
						Log.Print(err.Error() + " (" + line + ")")

						if strings.Contains(err.Error(), "runtime") {
							return
						}
					}

					go func() {
						time.Sleep(75 * time.Millisecond)
						ProcessLine(line) // Try the line again
					}()
				}
			}()
			handler.Impl(rawline)
		}
	}
}

// Main is the main loop.
func Main() {
	for {
		line, err := Conn.GetLine()
		if err != nil {
			break
		}

		debug("Got line")

		if Bursted {
			tasks <- line
		} else {
			debug("begin process line")
			ProcessLine(line)
			debug("End process line")
		}
	}

	wg.Wait()
}
