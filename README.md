Tetra
=====

Extended services in Go with Lua scripting

Tetra is more of a functional experiment than a services package right now. It 
still needs many things to be production ready, but here is what it has so far:

 - Yaml API
 - Lua script loading
 - Hooking on protocol events
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

You also need the following from moonrocks:

 - `lua_json`
 - `moonscript`
 - `yaml` (for runit-based autoconfiguration)

Run `./Tetra` in a tmux/dtach session. You can also run Tetra in Docker as long 
as you set `IRCD_PORT_6667_TCP_ADDR` to the address of your ircd. The Makefile 
in the project root automates this.
