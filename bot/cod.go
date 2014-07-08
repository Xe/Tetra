package tetra

import (
	"bufio"
	"errors"
	"github.com/Xe/Tetra/1459"
	"log"
	"net"
	"net/textproto"
	"os"
	"strconv"
)

type Clients struct {
	ByNick map[string]Client
	ByUID  map[string]Client
	Tetra  *Tetra
}

func (clients *Clients) AddClient(client Client) {
	clients.ByNick[client.Nick()] = client
	clients.ByUID[client.Uid()] = client
}

func (clients *Clients) DelClient(client Client) (err error) {
	return
}

type Tetra struct {
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

func NewTetra() (tetra *Tetra) {
	tetra = &Tetra{
		Conn: &Connection{
			Log: log.New(os.Stdout, "", log.LstdFlags),
		},
		Info: &Server{
			Name:  "tetra.int",
			Sid:   "420",
			Gecos: "Tetra in Go!",
		},
		Clients: &Clients{
			ByNick: make(map[string]Client),
			ByUID:  make(map[string]Client),
			Tetra:  tetra,
		},
		Channels: make(map[string]*Channel),
		Handlers: make(map[string]map[string]*Handler),
		Services: make(map[string]*ServiceClient),
		Servers:  make(map[string]*Server),
		Bursted:  false,
		nextuid:  100000,
	}

	tetra.AddHandler("EUID", func(line *r1459.RawLine) {
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
			tetra:   tetra,
		}

		tetra.Clients.AddClient(*client)
	})

	tetra.AddHandler("SJOIN", func(line *r1459.RawLine) {
		// :47G SJOIN 1404424869 #test +nt :@47GAAAABL
	})

	tetra.AddService("tetra", "Tetra", "user", "yolo-swag.com", "Tetra in Go!")

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

func (tetra *Tetra) AddService(service, nick, user, host, gecos string) (cli *ServiceClient) {
	cli = &ServiceClient{
		nick:    nick,
		user:    user,
		host:    host,
		VHost:   host,
		gecos:   gecos,
		account: "*",
		Ip:      "0",
		ts:      0,
		uid:     tetra.NextUID(),
	}

	tetra.Services[service] = cli

	tetra.Clients.AddClient(cli)

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
