package main

import (
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot"
)

func main() {
	bot := tetra.NewTetra("etc/config.json")

	bot.Connect("127.0.0.1", "6667")
	defer bot.Conn.Conn.Close()

	bot.Auth()

	for _, script := range bot.Config.Autoload {
		bot.LoadScript(script)
	}

	for _, sclient := range bot.Config.Services {
		bot.AddService(sclient.Name, sclient.Nick, sclient.User, sclient.Host, sclient.Gecos)
	}

	for _, client := range bot.Services {
		bot.Conn.SendLine(client.Euid())
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
				bot.Bursted = true
				if svc, ok := bot.Services["tetra"]; !ok {
					panic("No service bot!")
				} else {
					svc.Join("#services")
				}
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
