package main

import (
	"fmt"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot"
	"runtime"
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

	bot.AddCommand("tetra", "MEM",
		func(user *tetra.Client, message []string) string {
			stats := new(runtime.MemStats)
			runtime.ReadMemStats(stats)

			return fmt.Sprintf("Allocs: %d, Frees: %d, Bytes in use: %d, Scripts loaded: %d",
				stats.Mallocs, stats.Frees, stats.Alloc, len(bot.Scripts))
		})

	for _, client := range bot.Services {
		bot.Conn.SendLine(client.Euid())
		bot.Log.Printf("%#v", client)
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
