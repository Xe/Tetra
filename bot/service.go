package tetra

import (
	"errors"
	"time"
)

// AddService adds a new service Client to the network.
func (tetra *Tetra) AddService(service, nick, user, host, gecos, certfp string) (cli *Client) {
	cli = &Client{
		Nick:     nick,
		User:     user,
		Host:     "0",
		VHost:    host,
		Gecos:    gecos,
		Account:  nick,
		Ip:       "0",
		Ts:       time.Now().Unix(),
		Uid:      tetra.NextUID(),
		tetra:    tetra,
		Channels: make(map[string]*Channel),
		Server:   tetra.Info,
		Kind:     service,
		Commands: make(map[string]*Command),
		Certfp:   certfp,
	}

	tetra.Services[service] = cli

	tetra.Clients.AddClient(cli)

	if tetra.Bursted {
		tetra.Conn.SendLine(cli.Euid())
		if cli.Certfp != "" {
			tetra.Conn.SendLine(":%s ENCAP * CERTFP :%s", cli.Uid, cli.Certfp)
		}
	}

	tetra.Etcd.CreateDir("/tetra/scripts/" + cli.Kind, 0)

	return
}

// DelService deletes a service from the network or returns an error.
func (tetra *Tetra) DelService(service string) (err error) {
	if _, ok := tetra.Services[service]; !ok {
		return errors.New("No such service " + service)
	}

	client := tetra.Services[service]

	tetra.Clients.DelClient(client)
	client.Quit()

	return
}

