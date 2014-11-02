package tetra

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot/modes"
	"github.com/rcrowley/go-metrics"
)

func (tetra *Tetra) seedHandlers() {
	tetra.AddHandler("NICK", func(line *r1459.RawLine) {
		source := tetra.Clients.ByUID[line.Source]

		tetra.Clients.ChangeNick(source, line.Args[0])

		source.Nick = line.Args[0]
	})

	tetra.AddHandler("SQUIT", func(line *r1459.RawLine) {
		if line.Args[0] == tetra.Info.Sid {
			tetra.RunHook("SHUTDOWN")

			tetra.Log.Print("See you on the other side.")

			fmt.Println("Waiting for goroutines to settle... (5 seconds)")

			time.Sleep(5 * time.Second)

			os.Exit(0)
		}

		sid := line.Args[0]
		server := tetra.Servers[sid]

		// Remove all clients from the split server
		for uid, client := range tetra.Clients.ByUID {
			if strings.HasPrefix(uid, sid) {
				tetra.Clients.DelClient(client)
			}
		}

		delete(tetra.Servers, sid)

		for _, link := range server.Links {
			debugf("%#v", link)

			if link.Hops > server.Hops {
				for uid, client := range tetra.Clients.ByUID {
					if strings.HasPrefix(uid, link.Sid) {
						tetra.Clients.DelClient(client)
					}
				}

				delete(tetra.Servers, link.Sid)
			}
		}
	})

	tetra.AddHandler("ERROR", func(line *r1459.RawLine) {
		tetra.Log.Fatal(line.Raw)
	})

	tetra.AddHandler("PRIVMSG", func(line *r1459.RawLine) {
		source := tetra.Clients.ByUID[line.Source]
		destination := line.Args[0]
		message := strings.Split(line.Args[1], " ")[1:] // Don't repeat the verb

		if destination[0] == '#' {
			return
		} else {
			var ok bool
			_, ok = tetra.Clients.ByUID[destination]

			if !ok {
				tetra.Log.Fatal("got a message from a ghost client. We are out of sync.")
			}
		}

		client := tetra.Clients.ByUID[destination]
		verb := strings.ToUpper(strings.Split(line.Args[1], " ")[0])

		go func() {
			if command, ok := client.Commands[verb]; ok {
				if command.NeedsOper && !source.IsOper() {
					client.Notice(source, "Permission denied.")
					return
				}

				reply := command.Impl(source, client, message)

				if command.NeedsOper && reply != "" {
					client.ServicesLog(tetra.Clients.ByUID[source.Target()].Nick + ": " + command.Verb + ": " + reply)
				}

				client.Notice(source, reply)
			} else {
				client.Notice(source, "No such command "+verb)
			}
		}()
	})

	tetra.AddHandler("PRIVMSG", func(line *r1459.RawLine) {
		source := tetra.Clients.ByUID[line.Source]
		destination := line.Args[0]
		text := line.Args[1]

		if destination[0] != '#' {
			return
		}

		channel, ok := tetra.Channels[strings.ToUpper(destination)]
		if !ok {
			tetra.Log.Fatalf("Recieved CHANMSG from %s which is unknown. Panic.", destination)
		}

		if channel.Name == strings.ToUpper(tetra.Config.General.SnoopChan) {
			if strings.HasSuffix(source.Nick, "Serv") {
				tetra.RunHook(strings.ToUpper(source.Nick)+"-SERVICELOG", strings.Split(text, " "))
			}
		} else {
			for kind, client := range tetra.Services {
				if _, ok := client.Channels[channel.Target()]; ok {
					tetra.RunHook(strings.ToUpper(kind)+"-CHANMSG", source, channel, strings.Split(text, " "))
				}
			}
		}
	})

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
			Metadata: make(map[string]string),
		}

		tetra.Clients.AddClient(client)

		if tetra.Bursted {
			tetra.RunHook("NEWCLIENT", client)
		}
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
			Metadata: make(map[string]string),
		}

		client.Server.AddClient()

		tetra.Clients.AddClient(client)

		if tetra.Bursted {
			tetra.RunHook("NEWCLIENT", client)
		}
	})

	tetra.AddHandler("BMASK", func(line *r1459.RawLine) {
		// :42F BMASK 1414880311 #services b :fun!*@*
		channame := line.Args[1]
		bankind, ok := modes.CHANMODES[0][line.Args[2]]

		if !ok {
			return
		}

		masks := strings.Split(line.Args[3], " ")
		var channel *Channel

		if mychannel, ok := tetra.Channels[channame]; ok {
			channel = mychannel
		} else {
			return
		}

		channel.Lists[bankind] = append(channel.Lists[bankind], masks...)
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

			if tetra.Bursted {
				tetra.RunHook("JOINCHANNEL", cu)
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

	tetra.AddHandler("TMODE", func(line *r1459.RawLine) {
		channame := line.Args[1]
		modestring := line.Args[2]
		params := line.Args[3:]

		paramcounter := 0
		set := true

		channel, ok := tetra.Channels[strings.ToUpper(channame)]
		if !ok {
			return
		}

		for _, modechar := range modestring {
			mode := string(modechar)
			switch mode {
			case "+":
				set = true
			case "-":
				set = false
			default:
				if flag, ok := modes.CHANMODES[0][mode]; ok { // list-like mode
					if set {
						channel.Lists[flag] = append(channel.Lists[flag], params[paramcounter])
					} else {
						for i, str := range channel.Lists[flag] {
							if str == params[paramcounter] {
								channel.Lists[flag] = append(channel.Lists[flag][:i], channel.Lists[flag][i+1:]...)
							}
						}
					}
				} else if _, ok := modes.CHANMODES[1][mode]; ok { // channel set flag
					if set {
						channel.Modes = channel.Modes | modes.CHANMODES[1][mode]
					} else {
						channel.Modes = channel.Modes &^ (modes.CHANMODES[1][mode])
					}
				} else if _, ok := modes.CHANMODES[2][mode]; ok { // channel prefix flag
					if set {
						channel.Clients[params[paramcounter]].Prefix =
							channel.Clients[params[paramcounter]].Prefix | modes.CHANMODES[2][mode]
					} else {
						channel.Clients[params[paramcounter]].Prefix =
							channel.Clients[params[paramcounter]].Prefix &^ (modes.CHANMODES[2][mode])
					}
					paramcounter++
				} else { // modes that exist yet we don't care about
					if (mode == "j" || mode == "f") && set {
						paramcounter++
					} else if mode == "k" {
						paramcounter++
					}
				}
			}
		}
	})

	tetra.AddHandler("JOIN", func(line *r1459.RawLine) {
		client := tetra.Clients.ByUID[line.Source]
		channel := tetra.Channels[strings.ToUpper(line.Args[1])]

		channel.AddChanUser(client)

		tetra.RunHook("JOINCHANNEL", channel.Clients[client.Uid])
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

		tetra.RunHook("CLIENTQUIT", client)

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

		var err error
		server.Hops, err = strconv.Atoi(line.Args[1])
		if err != nil {
			return
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

	// Handle ENCAP by sending out a hook in the form of ENCAP-VERB.
	tetra.AddHandler("ENCAP", func(line *r1459.RawLine) {
		tetra.RunHook("ENCAP-"+line.Args[1], line.Source, line.Args[2:])
	})

	tetra.NewHook("ENCAP-GCAP", func(args ...interface{}) {
		if len(args) != 2 {
			return
		}

		var sid string
		var caps []string

		if targSid, ok := args[0].(string); ok {
			sid = targSid
		} else {
			return
		}

		if targCaps, ok := args[1].([]string); ok {
			caps = targCaps
		} else {
			return
		}

		server, ok := tetra.Servers[sid]
		if !ok {
			tetra.Log.Fatalf("Unknown server by ID %s. We are out of sync.", sid)
		}

		server.Capab = caps
	})

	tetra.NewHook("ENCAP-SU", func(parv ...interface{}) {
		var args []string
		var ok bool
		if args, ok = parv[1].([]string); !ok {
			return
		}

		if len(args) > 2 {
			return
		}

		var source *Client
		var accname string

		if source, ok = tetra.Clients.ByUID[args[0]]; !ok {
			return
		}

		if len(args) == 1 {
			accname = "*"
		} else {
			accname = args[1]
		}

		source.Account = accname
	})
}
