Setting up Tetra
================

#### Note

Tetra is only tested on Ubuntu Linux. Due to some specific cgo packages, it 
needs Ubuntu on the host machine.

Setup
-----

Tetra requires the following backing services:

- [etcd](http://github.com/coreos/etcd)
- [Influxdb](http://influxdb.com)
- [Elemental-IRCd](http://github.com/elemental-ircd/elemental-ircd)
- [Atheme Services](http://atheme.org)

Tetra also requires the following package from Moonrocks:

- `yaml`
- `moonscript`
- `json4lua`
- `luasocket`

The included Dockerfile will show what commands you would need to run on 
a fresh install of Ubuntu.

Configuration
-------------

Copy the default configuration to `etc/config.yaml` and open it in your 
favorite text editor.

The default configuration will have `tetra` and `chatbot` loaded. In order to 
load other services such as `mongblocker` or `zuul`, please add lines to that 
block that look something like this:

```
  - nick:   Zuul
    user:   guardian
    host:   yolo-swag.com
    gecos:  The gatekeeper
    name:   zuul
    certfp: 02438d05-b3e5-47f1-babd-ebd727169b8c
```

Linking
-------

In your elemental-ircd config, add a few lines similar to this:

```
connect "tetra.int" {
    host = "127.0.0.1";
    send_password = "shameless";
    accept_password = "shameless";
    port = 6667;
    class = "server;
};

service {
    name = "tetra.int";
};
```

It is critical that Tetra have the permissions of a services server so that it 
can do what it needs to do.

Then run Tetra inside either screen, tmux or dtach. Daemonization is not 
supported.
