// Command main starts the program.
package main

import (
	"github.com/Xe/Tetra/bot"
)

func main() {
	bot := tetra.NewTetra("etc/config.json")

	bot.Connect(bot.Config.Uplink.Host, bot.Config.Uplink.Port)
	defer bot.Conn.Conn.Close()

	bot.Auth()
	bot.StickConfig()
	bot.WebApp()

	for {
		line, err := bot.Conn.GetLine()
		if err != nil {
			panic(err)
		}

		bot.ProcessLine(line)
	}
}
