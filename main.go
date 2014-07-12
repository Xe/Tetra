package main

import (
	"fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot"
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
			if rawline.Source == "" {
				if !bot.Bursted {
					bot.Burst()
					bot.Log.Printf("Bursted!")
				}
				bot.Conn.SendLine("PONG :" + rawline.Args[0])
			} else {
				bot.Conn.SendLine(":%s PONG %s :%s", bot.Info.Sid, bot.Info.Name, rawline.Source)
			}
		}

		if _, present := bot.Handlers[rawline.Verb]; present {
			for _, handler := range bot.Handlers[rawline.Verb] {
				func() {
					defer func() {
						if r := recover(); r != nil {
							str := fmt.Sprintf("Recovered in handler %s (%s): %#v",
								handler.Verb, handler.Uuid, r)
							bot.Log.Printf(str)
							bot.Services["tetra"].ServicesLog(str)
						}
					}()
					if bot.Bursted {
						go handler.Impl(rawline)
					} else {
						handler.Impl(rawline)
					}
				}()
			}
		}
	}
}
