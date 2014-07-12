package tetra

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot/modes"
	"log"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"
)

type Clients struct {
	ByNick map[string]*Client
	ByUID  map[string]*Client
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
		},
		Info: &Server{},
		Clients: &Clients{
			ByNick: make(map[string]*Client),
			ByUID:  make(map[string]*Client),
			Tetra:  tetra,
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
		Uplink:   &Server{},
	}

	tetra.Info = &Server{
		Sid:   tetra.Config.Server.Sid,
		Name:  tetra.Config.Server.Name,
		Gecos: tetra.Config.Server.Gecos,
		Links: []*Server{tetra.Uplink},
	}

	go tetra.Conn.sendLinesWait()

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
	})

	tetra.AddHandler("SID", func(line *r1459.RawLine) {
		// <<< :42F SID cod.int 2 752 :Cod fishy
		parent := tetra.Servers[line.Source]

		server := &Server{
			Name:  line.Args[0],
			Gecos: line.Args[3],
			Sid:   line.Args[2],
			Links: []*Server{parent},
		}

		parent.Links = append(parent.Links, server)

		tetra.Servers[server.Sid] = server
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
	})

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

	tetra.Bursted = true
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
		Host:     host,
		VHost:    host,
		Gecos:    gecos,
		Account:  nick,
		Ip:       "0",
		Ts:       time.Now().Unix(),
		Uid:      tetra.NextUID(),
		tetra:    tetra,
		Channels: make(map[string]*Channel),
		Server:   tetra.Info,
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
