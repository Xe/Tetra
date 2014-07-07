package cod

import (
	"bufio"
	"errors"
	"github.com/cod-services/cod/1459"
	"log"
	"net"
	"net/textproto"
	"os"
	"strconv"
)

type Clients struct {
	ByNick map[string]Client
	ByUID  map[string]Client
	Cod    *Cod
}

func (clients *Clients) AddClient(client Client) {
	clients.ByNick[client.Nick()] = client
	clients.ByUID[client.Uid()] = client
}

func (clients *Clients) DelClient(client Client) (err error) {
	return
}

type Cod struct {
	Conn     *Connection
	Info     *Server
	Clients  *Clients
	Channels map[string]*Channel
	Bursted  bool
	Handlers map[string]map[string]*Handler
	Services map[string]*ServiceClient
	Servers  map[string]*Server
	nextuid  int
	//Config *Config
}

func NewCod() (cod *Cod) {
	cod = &Cod{
		Conn: &Connection{
			Log: log.New(os.Stdout, "", log.LstdFlags),
		},
		Info: &Server{
			Name:  "cod.int",
			Sid:   "420",
			Gecos: "Cod in Go!",
		},
		Clients: &Clients{
			ByNick: make(map[string]Client),
			ByUID:  make(map[string]Client),
			Cod:    cod,
		},
		Channels: make(map[string]*Channel),
		Handlers: make(map[string]map[string]*Handler),
		Services: make(map[string]*ServiceClient),
		Servers:  make(map[string]*Server),
		Bursted:  false,
		nextuid:  100000,
	}

	cod.AddHandler("EUID", func(line *r1459.RawLine) {
		// :47G EUID xena 1 1404369238 +ailoswxz xena staff.yolo-swag.com 0::1 47GAAAABK 0::1 * :Xena
		nick := line.Args[0]
		user := line.Args[4]
		host := line.Args[5]
		ip := line.Args[8]
		uid := line.Args[7]

		client := &RemoteClient{
			nick:    nick,
			user:    user,
			VHost:   host,
			host:    line.Args[6],
			uid:     uid,
			Ip:      ip,
			account: line.Args[9],
			gecos:   line.Args[10],
			cod:     cod,
		}

		cod.Clients.AddClient(*client)
	})

	cod.AddHandler("SJOIN", func(line *r1459.RawLine) {
		// :47G SJOIN 1404424869 #test +nt :@47GAAAABL
	})

	cod.AddService("cod", "Cod", "user", "yolo-swag.com", "Cod in Go!")

	return
}

func (cod *Cod) NextUID() string {
	cod.nextuid ++
	return cod.Info.Sid + strconv.Itoa(cod.nextuid)
}

func (cod *Cod) Connect(host, port string) (err error) {
	cod.Conn.Conn, err = net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}

	cod.Conn.Reader = bufio.NewReader(cod.Conn.Conn)
	cod.Conn.Tp = textproto.NewReader(cod.Conn.Reader)

	return
}

func (cod *Cod) AddService(service, nick, user, host, gecos string) (cli *ServiceClient) {
	cli = &ServiceClient{
		nick:  nick,
		user:  user,
		host:  host,
		VHost: host,
		gecos: gecos,
		account: "*",
		Ip: "0",
		ts: 0,
		uid: cod.NextUID(),
	}

	cod.Services[service] = cli

	cod.Clients.AddClient(cli)

	return
}

func (cod *Cod) DelService(service string) (err error) {
	if _, ok := cod.Services[service]; !ok {
		panic(errors.New("No such service " + service))
	}

	client := cod.Services[service]

	cod.Clients.DelClient(client)

	return
}

func (cod *Cod) GetConn() *net.Conn {
	return &cod.Conn.Conn
}
