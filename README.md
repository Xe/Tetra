# Tetra

[![GoDoc](https://godoc.org/github.com/Xe/Tetra?status.svg)](https://godoc.org/github.com/Xe/Tetra)

Tetra is an extended services package for TS6 IRC daemons with Lua and 
Moonscript support.

## Features

- JSON API
- Lua / Moonscript script loading
- Hooking on protocol events
- Hooking on arbitrary events
- Client/Channel/Server link tracking
- Statistics via influxdb
- Persistent data via etcd
- Atheme integration

### Things still in progress

- Feature parity with Cod
- Documentation on migration from Cod to Tetra
- Scripts being able to define webpages

Building a script for Tetra is as easy as:

```moonscript
Command "PING", ->
  "PONG"
```

## Installation

### From git

You need the following buildtime dependencies:

- `liblua5.1-dev`
- `golang`
- `libsqlite3-dev`

Example commands to set up the global environment needed for Tetra are in the 
included `Dockerfile`.

```console
$ go get github.com/Xe/Tetra
$ cd $GOPATH/github.com/Xe/Tetra
```

Continue with configuration.

### From a tarball

Install `liblua5.1-dev` then extract the tarball and continue with
configuration.

## Configuration

Look at the example config, copy it to `etc/config.yaml` or set
`TETRA_CONFIG_PATH` to a file on the disk. Edit the config to your needs.

## Running

You need to set up `etcd` for runtime key->value support for Tetra. You also 
need to set up InfluxxDB if you want to have Tetra track channel and server 
statistics. An instance of Atheme with the XMLRPC module loaded is required.

You need the following lua rocks:

- `luasocket`
- `moonscript`
- `yaml`
- `json4lua`
- `lsqlite3`

All are available in [moonrocks](http://rocks.moonscript.org).
