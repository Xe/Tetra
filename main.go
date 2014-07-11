package main

import (
	"fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot"
	"strings"
)

func main() {
	bot := tetra.NewTetra("etc/config.json")

	bot.Connect("127.0.0.1", "6667")
	defer bot.Conn.Conn.Close()

	bot.Auth()

	for _, sclient := range bot.Config.Services {
		bot.AddService(sclient.Name, sclient.Nick, sclient.User, sclient.Host, sclient.Gecos)
	}

	for _, script := range bot.Config.Autoload {
		bot.LoadScript(script)
	}

	for {
		line, err := bot.Conn.GetLine()
		if err != nil {
			panic(err)
		}

		rawline := r1459.NewRawLine(line)

		bot.Conn.Log.Printf("<<< %s", line)

		if rawline.Verb == "PING" {
			if !bot.Bursted {
				for _, client := range bot.Services {
					bot.Conn.SendLine(client.Euid())
					client.Join(bot.Config.Server.SnoopChan)
				}

				for _, channel := range bot.Channels {
					for uid, _ := range channel.Clients {
						if !strings.HasPrefix(uid, bot.Info.Sid) {
							continue
						}
						str := fmt.Sprintf(":%s SJOIN %d %s + :%s", bot.Info.Sid, channel.Ts, channel.Name, uid)
						bot.Conn.SendLine(str)
					}
				}

				bot.Bursted = true
			}
			bot.Conn.SendLine("PONG :" + rawline.Args[0])
		}

		if _, present := bot.Handlers[rawline.Verb]; present {
			for _, handler := range bot.Handlers[rawline.Verb] {
				if bot.Bursted {
					go handler.Impl(rawline)
				} else {
					handler.Impl(rawline)
				}
			}
		}
	}
}
