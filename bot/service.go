package tetra

import (
	"errors"
	"time"
)

// AddService adds a new service Client to the network.
func AddService(service, nick, user, host, gecos, certfp string) (cli *Client) {
	cli = &Client{
		Nick:     nick,
		User:     user,
		Host:     "0",
		VHost:    host,
		Gecos:    gecos,
		Account:  "*",
		Ip:       "0",
		Ts:       time.Now().Unix(),
		Uid:      NextUID(),
		Channels: make(map[string]*Channel),
		Server:   Info,
		Kind:     service,
		Commands: make(map[string]*Command),
		Certfp:   certfp,
		Metadata: make(map[string]string),
	}

	Services[service] = cli

	Clients.AddClient(cli)

	if Bursted {
		Conn.SendLine(cli.Euid())
		if cli.Certfp != "" {
			Conn.SendLine(":%s ENCAP * CERTFP :%s", cli.Uid, cli.Certfp)
		}
	}

	Etcd.CreateDir("/tetra/scripts/"+cli.Kind, 0)

	return
}

// DelService deletes a service from the network or returns an error.
func DelService(service string) (err error) {
	if _, ok := Services[service]; !ok {
		return errors.New("No such service " + service)
	}

	client := Services[service]

	Clients.DelClient(client)
	client.Quit()

	return
}
