package tetra

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/modes"
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
}

func NewTetra(cpath string) (tetra *Tetra) {
	config, err := NewConfig(cpath)
	if err != nil {
		fmt.Printf("No config file %s\n", cpath)
		panic(err)
	}

	tetra = &Tetra{
		Conn: &Connection{
			Log: log.New(os.Stdout, "CONN ", log.LstdFlags),
		},
		Info: &Server{
			Name:  "tetra.int",
			Sid:   "420",
			Gecos: "Tetra in Go!",
		},
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
	}

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
			Nick:    nick,
			User:    user,
			VHost:   host,
			Host:    line.Args[6],
			Uid:     uid,
			Ip:      ip,
			Account: line.Args[9],
			Gecos:   line.Args[10],
			tetra:   tetra,
			Umodes:  modeflags,
		}

		tetra.Clients.AddClient(client)
	})

	tetra.AddHandler("SJOIN", func(line *r1459.RawLine) {
		// :47G SJOIN 1404424869 #test +nt :@47GAAAABL

		ts := line.Args[0]
		name := line.Args[1]
		cmodes := line.Args[2]
		users := line.Args[3]

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

	tetra.AddHandler("CHGHOST", func(line *r1459.RawLine) {
		client := tetra.Clients.ByUID[line.Args[0]]
		client.VHost = line.Args[1]
	})

	tetra.AddHandler("QUIT", func(line *r1459.RawLine) {
		client := tetra.Clients.ByUID[line.Source]
		tetra.Clients.DelClient(client)
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

func (tetra *Tetra) AddService(service, nick, user, host, gecos string) (cli *Client) {
	cli = &Client{
		Nick:    nick,
		User:    user,
		Host:    host,
		VHost:   host,
		Gecos:   gecos,
		Account: nick,
		Ip:      "0",
		Ts:      time.Now().Unix(),
		Uid:     tetra.NextUID(),
		tetra:   tetra,
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
		panic(errors.New("No such service " + service))
	}

	client := tetra.Services[service]

	tetra.Clients.DelClient(client)

	return
}

func (tetra *Tetra) GetConn() *net.Conn {
	return &tetra.Conn.Conn
}
