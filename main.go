/*
Command Tetra is an extended services package for TS6 IRC daemons with Lua and
Moonscript support.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Xe/Tetra/bot"
)

var (
	config = flag.String("config", "etc/config.yaml", "configuration file for Tetra to use (or TETRA_CONFIG_PATH)")
)

func main() {
	flag.Parse()

	confloc := *config

	if envvar := os.Getenv("TETRA_CONFIG_PATH"); envvar != "" {
		confloc = envvar
	}

	if _, err := os.Open(confloc); err != nil {
		fmt.Fprintln(os.Stderr, "Please add your config at "+*config)
		os.Exit(1)
	}

	fmt.Printf("Using config file %s\n", confloc)

	cmd := exec.Command("moonc", ".")
	cmd.Dir = "./lib"

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	tetra.NewTetra(confloc)

	tetra.Connect(tetra.ActiveConfig.Uplink.Host, tetra.ActiveConfig.Uplink.Port)
	defer tetra.Conn.Conn.Close()

	tetra.Auth()
	tetra.StickConfig()
	go tetra.WebApp()

	tetra.Main()
}
