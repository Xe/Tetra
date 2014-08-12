# Tetra
--
Tetra is an extended services package for TS6 IRC daemons with Lua and 
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

Building a script for Tetra is as easy as:

```moonscript
require "modules/base"

command("PING") .. ->
  return "PONG"
```

## Installation

### From git

You need the following buildtime dependencies:

- `lua5.1`
- `golang`

```console
$ go get github.com/Xe/Tetra
$ cd $GOPATH/github.com/Xe/Tetra
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
- `moonscript`
- `etlua`

All are available in [moonrocks](http://rocks.moonscript.org).
