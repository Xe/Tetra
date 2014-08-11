/*
Command Tetra is an extended services package for TS6 IRC daemons with Lua and
Moonscript support.

Tetra is more of a functional experiment than a services package right now. It 
still needs many things to be production ready, but here is what it has so far:

 - Yaml API
 - Lua / Moonscript script loading
 - Hooking on protocol events
 - Hooking on arbitrary events
 - Client/Channel/Server link tracking
 - Statistics via influxdb

Things still in progress:

 - Feature parity with Cod
 - Documentation on migration from Cod to Tetra
 - Atheme integration
 - Scripts being able to define webpages

## Installation

### From git

You need the following buildtime dependencies:

 - `lua5.1`
 - `golang`

```console

$ go get github.com/Xe/Tetra

$ cd $GOPATH/

```

Continue with configuration.

### From a tarball

Install `liblua5.1-dev` then extract the tarball and continue with 
configuration.

## Configuration

Look at the example config, copy it to `etc/config.json` or set 
`TETRA_CONFIG_PATH` to a file on the disk. Edit the config to your needs.

## Running

You need the following lua rocks:

 - `luasocket`

*/
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
		if file, err = os.Open("etc/config.yaml"); err != nil {
			fmt.Fprintln(os.Stderr, "Please either set TETRA_CONFIG_PATH to the location of the configuration file or add your config at etc/config.yaml")
			os.Exit(1)
		} else {
			confloc = "etc/config.yaml"
			file.Close()
		}
	}

	fmt.Printf("Config file %s\n", confloc)

	bot := tetra.NewTetra(confloc)

	bot.Connect(bot.Config.Uplink.Host, bot.Config.Uplink.Port)
	defer bot.Conn.Conn.Close()

	bot.Auth()
	bot.StickConfig()
	bot.WebApp()

	bot.Main()
}
