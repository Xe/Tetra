package main

import (
	"github.com/Xe/Tetra/bot"
)

func main() {
	bot := tetra.NewTetra("etc/config.json")

	bot.Connect("127.0.0.1", "6667")
	defer bot.Conn.Conn.Close()

	bot.Auth()
	bot.StickConfig()

	for {
		line, err := bot.Conn.GetLine()
		if err != nil {
			panic(err)
		}

		bot.ProcessLine(line)
	}
}
