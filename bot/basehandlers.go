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

func seedHandlers() {
	AddHandler("NICK", func(line *r1459.RawLine) {
		source := Clients.ByUID[line.Source]

		Clients.ChangeNick(source, line.Args[0])

		source.Nick = line.Args[0]
	})

	AddHandler("SQUIT", func(line *r1459.RawLine) {
		if line.Args[0] == Info.Sid {
			RunHook("SHUTDOWN")

			Log.Print("See you on the other side.")

			fmt.Println("Waiting for goroutines to settle... (5 seconds)")

			time.Sleep(5 * time.Second)

			os.Exit(0)
		}

		sid := line.Args[0]
		server, ok := Servers[sid]
		if !ok {
			Log.Fatalf("Unknown server by ID %s", sid)
		}

		// Remove all clients from the split server
		for uid, client := range Clients.ByUID {
			if strings.HasPrefix(uid, sid) {
				Clients.DelClient(client)
			}
		}

		delete(Servers, sid)

		for _, link := range server.Links {
			debugf("%#v", link)

			if link.Hops > server.Hops {
				for uid, client := range Clients.ByUID {
					if strings.HasPrefix(uid, link.Sid) {
						Clients.DelClient(client)
					}
				}

				delete(Servers, link.Sid)
			}
		}
	})

	AddHandler("ERROR", func(line *r1459.RawLine) {
		Log.Fatal(line.Raw)
	})

	AddHandler("PRIVMSG", func(line *r1459.RawLine) {
		source, ok := Clients.ByUID[line.Source]
		if !ok {
			panic(fmt.Errorf("Cannot find client by UID %s", line.Source))
		}

		destination := line.Args[0]
		message := strings.Split(line.Args[1], " ")[1:] // Don't repeat the verb

		if destination[0] == '#' {
			return
		} else {
			var ok bool
			_, ok = Clients.ByUID[destination]

			if !ok {
				Log.Fatal("got a message from a ghost client. We are out of sync.")
			}
		}

		if line.Args[1][0] == '\x01' {
			return
		}

		client := Clients.ByUID[destination]
		verb := strings.ToUpper(strings.Split(line.Args[1], " ")[0])

		go func() {
			if command, ok := client.Commands[verb]; ok {
				if command.NeedsOper && !source.IsOper() {
					client.Notice(source, "Permission denied.")
					return
				}

				reply := command.Impl(source, client, message)

				if command.NeedsOper && reply != "" {
					client.ServicesLog(Clients.ByUID[source.Target()].Nick + ": " + command.Verb + ": " + reply)
				}

				client.Notice(source, reply)
			} else {
				client.Notice(source, "No such command "+verb)
			}
		}()
	})

	AddHandler("PRIVMSG", func(line *r1459.RawLine) {
		source, ok := Clients.ByUID[line.Source]
		if !ok {
			panic(fmt.Errorf("Cannot find client by UID %s", line.Source))
		}

		destination := line.Args[0]
		text := line.Args[1]

		if destination[0] != '#' {
			return
		}

		channel, ok := Channels[strings.ToUpper(destination)]
		if !ok {
			Log.Fatalf("Recieved CHANMSG from %s which is unknown. Panic.", destination)
		}

		if strings.ToUpper(channel.Name) == strings.ToUpper(ActiveConfig.General.SnoopChan) {
			if strings.HasSuffix(source.Nick, "Serv") {
				RunHook(strings.ToUpper(source.Nick)+"-SERVICELOG", strings.Split(text, " "))
			}
		} else {
			for kind, client := range Services {
				if _, ok := client.Channels[channel.Target()]; ok {
					RunHook(strings.ToUpper(kind)+"-CHANMSG", source, channel, strings.Split(text, " "))
				}
			}
		}
	})

	AddHandler("PRIVMSG", func(line *r1459.RawLine) {
		if line.Args[0][0] == '#' {
			return
		}

		if line.Args[1][0] != '\x01' {
			return
		}

		source, ok := Clients.ByUID[line.Source]
		if !ok {
			panic(fmt.Errorf("Cannot find client by UID %s", line.Source))
		}

		destination := Clients.ByUID[line.Args[0]]
		text := line.Args[1]

		verb := strings.Split(text, " ")[0]
		verb = verb[1 : len(verb)-1]

		switch verb {
		case "VERSION":
			destination.Notice(source, "\x01VERSION Tetra\x01")
		default:
			destination.Notice(source, "Unknown CTCP "+verb)
		}
	})

	AddHandler("UID", func(line *r1459.RawLine) {
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
			Umodes:   modeflags,
			Channels: make(map[string]*Channel),
			Server:   Servers[line.Source],
			Metadata: make(map[string]string),
		}

		Clients.AddClient(client)

		if Bursted {
			RunHook("NEWCLIENT", client)
		}
	})

	AddHandler("EUID", func(line *r1459.RawLine) {
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
			Umodes:   modeflags,
			Channels: make(map[string]*Channel),
			Server:   Servers[line.Source],
			Metadata: make(map[string]string),
		}

		client.Server.AddClient()

		Clients.AddClient(client)

		if Bursted {
			RunHook("NEWCLIENT", client)
		}
	})

	AddHandler("BMASK", func(line *r1459.RawLine) {
		// :42F BMASK 1414880311 #services b :fun!*@*
		channame := line.Args[1]
		bankind, ok := modes.CHANMODES[0][line.Args[2]]
		if !ok {
			panic(fmt.Errorf("Unknown channel %s", line.Args[2]))
		}

		masks := strings.Split(line.Args[3], " ")
		var channel *Channel

		if mychannel, ok := Channels[channame]; ok {
			channel = mychannel
		} else {
			panic(fmt.Errorf("Unknown channel %s", line.Args[2]))
		}

		channel.Lists[bankind] = append(channel.Lists[bankind], masks...)
	})

	AddHandler("SJOIN", func(line *r1459.RawLine) {
		// :47G SJOIN 1404424869 #test +nt :@47GAAAABL
		ts := line.Args[0]
		name := line.Args[1]
		cmodes := line.Args[2]

		if line.Raw[len(line.Raw)-1] == ':' {
			return
		}

		users := line.Args[len(line.Args)-1]

		var channel *Channel

		if mychannel, ok := Channels[name]; ok {
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

			channel = NewChannel(name, numberts)
			channel.Modes = modeflags
		}

		for _, user := range strings.Split(users, " ") {
			var uid string
			length := len(user)
			pfxcount := length - 9

			uid = user[pfxcount:]
			prefixes := user[:pfxcount]

			client := Clients.ByUID[uid]

			cu := channel.AddChanUser(client)

			for _, char := range prefixes {
				if _, ok := modes.PREFIXES[string(char)]; ok {
					cu.Prefix = modes.PREFIXES[string(char)] | cu.Prefix
				}
			}

			if Bursted {
				RunHook("JOINCHANNEL", cu)
			}
		}
	})

	AddHandler("MODE", func(line *r1459.RawLine) {
		var give bool = true
		client := Clients.ByUID[line.Args[0]]
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

	AddHandler("TMODE", func(line *r1459.RawLine) {
		channame := line.Args[1]
		modestring := line.Args[2]
		params := line.Args[3:]

		paramcounter := 0
		set := true

		channel, ok := Channels[strings.ToUpper(channame)]
		if !ok {
			panic(fmt.Errorf("Unknown channel name %s", channame))
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

					RunHook("LISTMODE", mode, flag, channel, params[paramcounter])
					paramcounter++
				} else if _, ok := modes.CHANMODES[1][mode]; ok { // channel set flag
					if set {
						channel.Modes = channel.Modes | modes.CHANMODES[1][mode]
					} else {
						channel.Modes = channel.Modes &^ (modes.CHANMODES[1][mode])
					}

					RunHook("SETMODE", mode, flag, channel)
				} else if _, ok := modes.CHANMODES[2][mode]; ok { // channel prefix flag
					if set {
						channel.Clients[params[paramcounter]].Prefix =
							channel.Clients[params[paramcounter]].Prefix | modes.CHANMODES[2][mode]
					} else {
						channel.Clients[params[paramcounter]].Prefix =
							channel.Clients[params[paramcounter]].Prefix &^ (modes.CHANMODES[2][mode])
					}

					RunHook("PREFIXMODE", mode, flag, channel, channel.Clients[params[paramcounter]])
					paramcounter++
				} else { // modes that exist yet we don't care about
					if (mode == "j" || mode == "f") && set {
						RunHook("PARAMETRICMODE", mode, flag, channel, params[paramcounter])
						paramcounter++
					} else if mode == "k" {
						RunHook("KEYMODE", mode, flag, channel, params[paramcounter])
						paramcounter++
					}
				}
			}
		}
	})

	AddHandler("JOIN", func(line *r1459.RawLine) {
		client, ok := Clients.ByUID[line.Source]
		if !ok {
			panic(fmt.Errorf("Unknown client %s", line.Source))
		}

		channel, ok := Channels[strings.ToUpper(line.Args[1])]
		if !ok {
			Log.Fatalf("Unknown channel %s, desync", strings.ToUpper(line.Args[1]))
		}

		cu := channel.AddChanUser(client)

		RunHook("JOINCHANNEL", cu)
	})

	AddHandler("PART", func(line *r1459.RawLine) {
		// <<< :42FAAAAAB PART #help
		channelname := strings.ToUpper(line.Args[0])
		client, ok := Clients.ByUID[line.Source]
		if !ok {
			Log.Fatalf("Unknown client %s, desync", line.Source)
		}

		channel, ok := Channels[channelname]
		if !ok {
			Log.Fatalf("Unknown channel %s, desync", channelname)
		}

		channel.DelChanUser(client)
	})

	AddHandler("KICK", func(line *r1459.RawLine) {
		// <<< :42FAAAAAB KICK #help 42FAAAAAB :foo
		channelname := strings.ToUpper(line.Args[0])
		client, ok := Clients.ByUID[line.Source]
		if !ok {
			panic(fmt.Errorf("Unknown client %s", line.Source))
		}

		channel, ok := Channels[channelname]
		if !ok {
			panic(fmt.Errorf("Unknown channel %s", line.Source))
		}

		channel.DelChanUser(client)
	})

	AddHandler("CHGHOST", func(line *r1459.RawLine) {
		client, ok := Clients.ByUID[line.Args[0]]
		if !ok {
			Log.Fatalf("Unknown client %s, desync", line.Source)
		}

		client.VHost = line.Args[1]
	})

	AddHandler("QUIT", func(line *r1459.RawLine) {
		client, ok := Clients.ByUID[line.Source]
		if !ok {
			Log.Fatalf("Unknown client %s, desync", line.Source)
		}

		RunHook("CLIENTQUIT", client)

		Clients.DelClient(client)

		for _, channel := range client.Channels {
			channel.DelChanUser(client)
		}

		client.Server.DelClient()
	})

	AddHandler("SID", func(line *r1459.RawLine) {
		// <<< :42F SID cod.int 2 752 :Cod fishy
		parent, ok := Servers[line.Source]
		if !ok {
			Log.Fatal("No server by ID " + line.Source + ", desync")
		}

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

		Servers[server.Sid] = server

		metrics.Register(server.Name+"_clients", server.Counter)
	})

	AddHandler("PASS", func(line *r1459.RawLine) {
		// <<< PASS shameless TS 6 :42F
		Uplink.Sid = line.Args[3]
		Servers[line.Args[3]] = Uplink
	})

	AddHandler("SERVER", func(line *r1459.RawLine) {
		// <<< SERVER fluttershy.yolo-swag.com 1 :shadowircd test server
		Uplink.Name = line.Args[0]
		Uplink.Gecos = line.Args[2]

		metrics.Register(Uplink.Name+"_clients", Uplink.Counter)
	})

	AddHandler("WHOIS", func(line *r1459.RawLine) {
		/*
			<<< :649AAAABQ WHOIS 376100000 :ShadowNET
			>>> :376 311 649AAAABQ ShadowNET fishie cod.services * :Cod IRC services
			>>> :376 312 649AAAABQ ShadowNET ardreth.shadownet.int :Cod IRC services
			>>> :376 313 649AAAABQ ShadowNET :is a Network Service
			>>> :376 318 649AAAABQ ShadowNET :End of /WHOIS list.
		*/

		target := line.Args[0]
		client := Clients.ByUID[target]
		source := Clients.ByUID[line.Source]

		if client.Kind == "" {
			return
		}

		temp := []string{
			fmt.Sprintf(":%s 311 %s %s %s %s * :%s", Info.Sid, source.Uid,
				client.Nick, client.User, client.VHost, client.Gecos),
			fmt.Sprintf(":%s 312 %s %s %s :%s", Info.Sid, source.Uid,
				client.Nick, Info.Name, Info.Gecos),
			fmt.Sprintf(":%s 313 %s %s :is a Network Service (%s)",
				Info.Sid, source.Uid, client.Nick, client.Kind),
			fmt.Sprintf(":%s 318 %s %s :End of /WHOIS list.", Info.Sid,
				source.Uid, client.Nick),
		}

		for _, line := range temp {
			Conn.SendLine(line)
		}
	})

	// Handle ENCAP by sending out a hook in the form of ENCAP-VERB.
	AddHandler("ENCAP", func(line *r1459.RawLine) {
		RunHook("ENCAP-"+line.Args[1], line.Source, line.Args[2:])
	})

	NewHook("ENCAP-GCAP", func(args ...interface{}) {
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

		server, ok := Servers[sid]
		if !ok {
			Log.Fatalf("Unknown server by ID %s. We are out of sync.", sid)
		}

		server.Capab = caps
	})

	NewHook("ENCAP-SU", func(parv ...interface{}) {
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

		if source, ok = Clients.ByUID[args[0]]; !ok {
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
