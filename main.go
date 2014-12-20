/*
Command Tetra is an extended services package for TS6 IRC daemons with Lua and
Moonscript support.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/Xe/Tetra/bot"
)

var (
	config = flag.String("config", "etc/config.yaml", "configuration file for Tetra to use (or TETRA_CONFIG_PATH)")
	procs  = flag.Int("procs", 1, "value for runtime.GOMAXPROCS")
)

func main() {
	flag.Parse()

	confloc := *config
	runtime.GOMAXPROCS(*procs)

	if envvar := os.Getenv("TETRA_CONFIG_PATH"); envvar != "" {
		confloc = envvar
	}

	if _, err := os.Open(confloc); err != nil {
		fmt.Fprintln(os.Stderr, "Please add your config at "+*config)
		os.Exit(1)
	}

	fmt.Printf("Using config file %s\n", confloc)

	bot := tetra.NewTetra(confloc)

	bot.Connect(bot.Config.Uplink.Host, bot.Config.Uplink.Port)
	defer bot.Conn.Conn.Close()

	bot.Auth()
	bot.StickConfig()
	go bot.WebApp()

	bot.Main()
}
