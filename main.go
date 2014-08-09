// Command main starts the program.
package main

import (
	"fmt"
	"os"

	"github.com/Xe/Tetra/bot"
)

func main() {
	confloc := os.Getenv("TETRA_CONFIG_PATH")

	if confloc == "" { // No user set config location
		var file *os.File
		var err error
		if file, err = os.Open("etc/config.yaml") ; err != nil {
			fmt.Fprintln(os.Stderr, "Please either set TETRA_CONFIG_PATH to the location of the configuration file or add your config at etc/config.yaml")
			os.Exit(1)
		} else {
			confloc = "etc/config.yaml"
			file.Close()
		}
	}

	bot := tetra.NewTetra(confloc)

	bot.Connect(bot.Config.Uplink.Host, bot.Config.Uplink.Port)
	defer bot.Conn.Conn.Close()

	bot.Auth()
	bot.StickConfig()
	bot.WebApp()

	bot.Main()
}
