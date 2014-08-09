package tetra

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot/modes"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/influxdb"
	"log"
	"log/syslog"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Clients struct {
	ByNick map[string]*Client
	ByUID  map[string]*Client
	Gauge  metrics.Gauge
	Tetra  *Tetra
}

func (clients *Clients) AddClient(client *Client) {
	clients.ByNick[client.Nick] = client
	clients.ByUID[client.Uid] = client
}

func (clients *Clients) DelClient(client *Client) (err error) {
	delete(clients.ByNick, client.Nick)
	delete(clients.ByUID, client.Uid)

	return
}

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
	nextuid  int
	Config   *Config
	Log      *log.Logger
	Uplink   *Server
	tasks    chan string
	wg       *sync.WaitGroup
}

func NewTetra(cpath string) (tetra *Tetra) {
	config, err := NewConfig(cpath)
	if err != nil {
		fmt.Printf("No config file %s\n", cpath)
		panic(err)
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
		Bursted:  false,
		nextuid:  100000,
		Config:   config,
		Log:      log.New(os.Stdout, "BOT ", log.LstdFlags),
		Uplink: &Server{
			Counter: metrics.NewGauge(),
		},
		tasks: make(chan string, 100),
		wg:    &sync.WaitGroup{},
	}

	tetra.Info = &Server{
		Sid:     tetra.Config.Server.Sid,
		Name:    tetra.Config.Server.Name,
		Gecos:   tetra.Config.Server.Gecos,
		Links:   []*Server{tetra.Uplink},
		Counter: nil,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			_ = sig
			for _, client := range tetra.Services {
				tetra.Conn.SendLine(":%s QUIT :Shutting down.", client.Uid)
			}
			tetra.Conn.SendLine("SQUIT %s :Bye", tetra.Info.Sid)
			tetra.Conn.Close()

			tetra.AddHandler("SQUIT", func(line *r1459.RawLine) {
				os.Exit(0)
			})
		}
	}()

	metrics.Register(tetra.Config.Server.Name+"_clients", tetra.Info.Counter)

	go tetra.Conn.sendLinesWait()

	tetra.AddHandler("UID", func(line *r1459.RawLine) {
		// <<< :0RS UID RServ 2 0 +Z rserv rserv.yolo-swag.com 0 0RSSR0001 :Ruby Services
		nick := line.Args[0]
		umodes := line.Args[3]
		user := line.Args[4]
		host := line.Args[5]
		ip := line.Args[6]
		uid := line.Args[7]

		// TODO: make this its own function somewhere?
		modeflags := 0

		for _, char := range umodes {
			if _, ok := modes.UMODES[string(char)]; ok {
				modeflags = modeflags | modes.UMODES[string(char)]
			}
		}

		client := &Client{
			Nick:     nick,
			User:     user,
			VHost:    host,
			Host:     line.Args[6],
			Uid:      uid,
			Ip:       ip,
			Account:  "*",
			Gecos:    line.Args[8],
			tetra:    tetra,
			Umodes:   modeflags,
			Channels: make(map[string]*Channel),
			Server:   tetra.Servers[line.Source],
		}

		tetra.Clients.AddClient(client)
	})

	tetra.AddHandler("EUID", func(line *r1459.RawLine) {
		// :47G EUID xena 1 1404369238 +ailoswxz xena staff.yolo-swag.com 0::1 47GAAAABK 0::1 * :Xena
		nick := line.Args[0]
		umodes := line.Args[3]
		user := line.Args[4]
		host := line.Args[5]
		ip := line.Args[8]
		uid := line.Args[7]

		// TODO: make this its own function somewhere?
		modeflags := 0

		for _, char := range umodes {
			if _, ok := modes.UMODES[string(char)]; ok {
				modeflags = modeflags | modes.UMODES[string(char)]
			}
		}

		client := &Client{
			Nick:     nick,
			User:     user,
			VHost:    host,
			Host:     line.Args[6],
			Uid:      uid,
			Ip:       ip,
			Account:  line.Args[9],
			Gecos:    line.Args[10],
			tetra:    tetra,
			Umodes:   modeflags,
			Channels: make(map[string]*Channel),
			Server:   tetra.Servers[line.Source],
		}

		client.Server.AddClient()

		tetra.Clients.AddClient(client)
	})

	tetra.AddHandler("SJOIN", func(line *r1459.RawLine) {
		// :47G SJOIN 1404424869 #test +nt :@47GAAAABL
		ts := line.Args[0]
		name := line.Args[1]
		cmodes := line.Args[2]

		if line.Raw[len(line.Raw)-1] == ':' {
			return
		}

		users := line.Args[len(line.Args)-1]

		var channel *Channel

		if mychannel, ok := tetra.Channels[name]; ok {
			channel = mychannel
		} else {
			// The ircd should never give an invalid TS.
			numberts, _ := strconv.ParseInt(ts, 10, 64)

			// TODO: make this its own function somewhere?
			modeflags := 0

			for _, char := range cmodes {
				if _, ok := modes.CHANMODES[1][string(char)]; ok {
					modeflags = modeflags | modes.CHANMODES[1][string(char)]
				}
			}

			channel = tetra.NewChannel(name, numberts)
			channel.Modes = modeflags
		}

		for _, user := range strings.Split(users, " ") {
			var uid string
			length := len(user)
			pfxcount := length - 9

			uid = user[pfxcount:]
			prefixes := user[:pfxcount]

			client := tetra.Clients.ByUID[uid]

			cu := channel.AddChanUser(client)

			for _, char := range prefixes {
				if _, ok := modes.PREFIXES[string(char)]; ok {
					cu.Prefix = modes.PREFIXES[string(char)] | cu.Prefix
				}
			}
		}
	})

	tetra.AddHandler("MODE", func(line *r1459.RawLine) {
		var give bool = true
		client := tetra.Clients.ByUID[line.Args[0]]
		modeflags := client.Umodes

		umodes := line.Args[1]

		for _, char := range umodes {
			if char == '+' {
				give = true
			} else if char == '-' {
				give = false
			}
			if _, ok := modes.UMODES[string(char)]; ok {
				if give {
					modeflags = modeflags | modes.UMODES[string(char)]
				} else {
					modeflags = modeflags & ^(modes.UMODES[string(char)])
				}
			}
		}

		client.Umodes = modeflags
	})

	tetra.AddHandler("JOIN", func(line *r1459.RawLine) {
		client := tetra.Clients.ByUID[line.Source]
		channel := tetra.Channels[strings.ToUpper(line.Args[1])]

		channel.AddChanUser(client)
	})

	tetra.AddHandler("PART", func(line *r1459.RawLine) {
		// <<< :42FAAAAAB PART #help
		channelname := strings.ToUpper(line.Args[0])
		channel := tetra.Channels[channelname]
		client := tetra.Clients.ByUID[line.Source]

		channel.DelChanUser(client)
	})

	tetra.AddHandler("KICK", func(line *r1459.RawLine) {
		// <<< :42FAAAAAB KICK #help 42FAAAAAB :foo
		channelname := strings.ToUpper(line.Args[0])
		channel := tetra.Channels[channelname]
		client := tetra.Clients.ByUID[line.Source]

		channel.DelChanUser(client)
	})

	tetra.AddHandler("CHGHOST", func(line *r1459.RawLine) {
		client := tetra.Clients.ByUID[line.Args[0]]
		client.VHost = line.Args[1]
	})

	tetra.AddHandler("QUIT", func(line *r1459.RawLine) {
		client := tetra.Clients.ByUID[line.Source]
		tetra.Clients.DelClient(client)

		for _, channel := range client.Channels {
			channel.DelChanUser(client)
		}

		client.Server.DelClient()
	})

	tetra.AddHandler("SID", func(line *r1459.RawLine) {
		// <<< :42F SID cod.int 2 752 :Cod fishy
		parent := tetra.Servers[line.Source]

		server := &Server{
			Name:    line.Args[0],
			Gecos:   line.Args[3],
			Sid:     line.Args[2],
			Links:   []*Server{parent},
			Counter: metrics.NewGauge(),
		}

		parent.Links = append(parent.Links, server)

		tetra.Servers[server.Sid] = server

		metrics.Register(server.Name+"_clients", server.Counter)
	})

	tetra.AddHandler("PASS", func(line *r1459.RawLine) {
		// <<< PASS shameless TS 6 :42F
		tetra.Uplink.Sid = line.Args[3]
		tetra.Servers[line.Args[3]] = tetra.Uplink
	})

	tetra.AddHandler("SERVER", func(line *r1459.RawLine) {
		// <<< SERVER fluttershy.yolo-swag.com 1 :shadowircd test server
		tetra.Uplink.Name = line.Args[0]
		tetra.Uplink.Gecos = line.Args[2]

		metrics.Register(tetra.Uplink.Name+"_clients", tetra.Uplink.Counter)
	})

	tetra.AddHandler("WHOIS", func(line *r1459.RawLine) {
		/*
			<<< :649AAAABQ WHOIS 376100000 :ShadowNET
			>>> :376 311 649AAAABQ ShadowNET fishie cod.services * :Cod IRC services
			>>> :376 312 649AAAABQ ShadowNET ardreth.shadownet.int :Cod IRC services
			>>> :376 313 649AAAABQ ShadowNET :is a Network Service
			>>> :376 318 649AAAABQ ShadowNET :End of /WHOIS list.
		*/

		target := line.Args[0]
		client := tetra.Clients.ByUID[target]
		source := tetra.Clients.ByUID[line.Source]

		temp := []string{
			fmt.Sprintf(":%s 311 %s %s %s %s * :%s", tetra.Info.Sid, source.Uid,
				client.Nick, client.User, client.VHost, client.Gecos),
			fmt.Sprintf(":%s 312 %s %s %s :%s", tetra.Info.Sid, source.Uid,
				client.Nick, tetra.Info.Name, tetra.Info.Gecos),
			fmt.Sprintf(":%s 313 %s %s :is a Network Service (%s)",
				tetra.Info.Sid, source.Uid, client.Nick, client.Kind),
			fmt.Sprintf(":%s 318 %s %s :End of /WHOIS list.", tetra.Info.Sid,
				source.Uid, client.Nick),
		}

		for _, line := range temp {
			tetra.Conn.SendLine(line)
		}
	})

	for i := 0; i < 16; i++ {
		tetra.wg.Add(1)
		go func() {
			for line := range tetra.tasks {
				tetra.ProcessLine(line)
			}
			tetra.wg.Done()
		}()
	}

	return
}

func (tetra *Tetra) NextUID() string {
	tetra.nextuid++
	return tetra.Info.Sid + strconv.Itoa(tetra.nextuid)
}

func (tetra *Tetra) Connect(host, port string) (err error) {
	tetra.Conn.Conn, err = net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}

	tetra.Conn.Reader = bufio.NewReader(tetra.Conn.Conn)
	tetra.Conn.Tp = textproto.NewReader(tetra.Conn.Reader)

	return
}

func (tetra *Tetra) Auth() {
	tetra.Conn.SendLine("PASS " + tetra.Config.Uplink.Password + " TS 6 :" + tetra.Config.Server.Sid)
	tetra.Conn.SendLine("CAPAB :QS EX IE KLN UNKLN ENCAP SERVICES EUID EOPMO")
	tetra.Conn.SendLine("SERVER " + tetra.Config.Server.Name + " 1 :" + tetra.Config.Server.Gecos)
}

func (tetra *Tetra) Burst() {
	for _, client := range tetra.Services {
		tetra.Conn.SendLine(client.Euid())
		client.Join(tetra.Config.Server.SnoopChan)
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

	w, _ := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
	go metrics.Syslog(metrics.DefaultRegistry, 60e9, w)

	metrics.Register("clientcount", tetra.Clients.Gauge)

	go tetra.GetNetworkStats()
	go tetra.GetChannelStats()

	go influxdb.Influxdb(metrics.DefaultRegistry, 5*time.Minute, &influxdb.Config{
		Host:     tetra.Config.Stats.Host,
		Database: tetra.Config.Stats.Database,
		Username: tetra.Config.Stats.Username,
		Password: tetra.Config.Stats.Password,
	})

	tetra.Bursted = true
}

func (tetra *Tetra) StickConfig() {
	for _, sclient := range tetra.Config.Services {
		tetra.AddService(sclient.Name, sclient.Nick, sclient.User, sclient.Host, sclient.Gecos)
	}

	for _, script := range tetra.Config.Autoload {
		tetra.LoadScript(script)
	}

}

func (tetra *Tetra) Quit() {
	for _, service := range tetra.Services {
		tetra.DelService(service.Kind)
		service.Quit()
	}

	tetra.Conn.SendLine("SQUIT :Goodbye!")
	tetra.Conn.Conn.Close()
}

func (tetra *Tetra) AddService(service, nick, user, host, gecos string) (cli *Client) {
	cli = &Client{
		Nick:     nick,
		User:     user,
		Host:     "0",
		VHost:    host,
		Gecos:    gecos,
		Account:  nick,
		Ip:       "0",
		Ts:       time.Now().Unix(),
		Uid:      tetra.NextUID(),
		tetra:    tetra,
		Channels: make(map[string]*Channel),
		Server:   tetra.Info,
		Kind:     service,
	}

	tetra.Services[service] = cli

	tetra.Clients.AddClient(cli)

	if tetra.Bursted {
		tetra.Conn.SendLine(cli.Euid())
	}

	return
}

func (tetra *Tetra) DelService(service string) (err error) {
	if _, ok := tetra.Services[service]; !ok {
		return errors.New("No such service " + service)
	}

	client := tetra.Services[service]

	tetra.Clients.DelClient(client)
	client.Quit()

	return
}

func (tetra *Tetra) GetConn() *net.Conn {
	return &tetra.Conn.Conn
}

func (tetra *Tetra) ProcessLine(line string) {
	rawline := r1459.NewRawLine(line)

	tetra.Conn.Log.Printf("<<< %s", line)

	defer func() {
		if r := recover(); r != nil {
			tetra.Conn.Log.Printf("<<< %s", line)
			str := fmt.Sprintf("Recovered from verb %s", rawline.Verb)
			tetra.Log.Printf(str)
			tetra.Services["tetra"].ServicesLog(str)
		}
	}()

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
			func() {
				defer func() {
					if r := recover(); r != nil {
						str := fmt.Sprintf("Recovered in handler %s (%s): %#v",
							handler.Verb, handler.Uuid, r)
						tetra.Log.Printf(str)
						tetra.Services["tetra"].ServicesLog(str)
					}
				}()
				if tetra.Bursted {
					go handler.Impl(rawline)
				} else {
					handler.Impl(rawline)
				}
			}()
		}
	}
}

func (t *Tetra) Main() {
	for {
		line, err := t.Conn.GetLine()
		if err != nil {
			break
		}

		t.tasks <- line
	}

	t.wg.Wait()
}
