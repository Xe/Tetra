package tetra

import (
	"errors"
	"github.com/rcrowley/go-metrics"
	"strings"
)

type ChanUser struct {
	Client  *Client
	Channel *Channel
	Prefix  int
}

// Implements Targeter
type Channel struct {
	Name    string
	Ts      int64
	Modes   int
	Clients map[string]*ChanUser
	Lists   map[int][]string
	Gauge   metrics.Gauge
}

func (tetra *Tetra) NewChannel(name string, ts int64) (c *Channel) {
	c = &Channel{
		Name:    name,
		Ts:      ts,
		Lists:   make(map[int][]string),
		Clients: make(map[string]*ChanUser),
		Modes:   0,
		Gauge:   metrics.NewGauge(),
	}

	tetra.Channels[c.Target()] = c

	metrics.Register(strings.ToUpper(name) + "_stats", c.Gauge)

	return
}

func (c *Channel) AddChanUser(client *Client) (cu *ChanUser) {

	cu = &ChanUser{
		Client:  client,
		Channel: c,
		Prefix:  0,
	}

	c.Clients[client.Uid] = cu

	client.Channels[c.Target()] = c

	return
}

func (c *Channel) DelChanUser(client *Client) (err error) {
	if _, ok := c.Clients[client.Uid]; !ok {
		return errors.New("Tried to delete nonexistent chanuser with uid " + client.Uid + " from " + c.Name)
	}

	delete(c.Clients, client.Uid)
	delete(client.Channels, c.Name)

	return nil
}

func (c *Channel) Target() string {
	return strings.ToUpper(c.Name)
}

func (c *Channel) IsChannel() bool {
	return true
}
