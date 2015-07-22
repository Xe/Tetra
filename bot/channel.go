package tetra

import (
	"errors"
	"strings"
	"sync"

	"github.com/rcrowley/go-metrics"
)

// Struct ChanUser is a wrapper around a Channel and a Client to represent membership
// in a Channel.
type ChanUser struct {
	Client  *Client
	Channel *Channel
	Prefix  int
}

// Struct Channel holds all the relevant data for an IRC channel. A lot of this
// is not just things defined in RFC 1459, but extensions like the TS.
// This implements Targeter
type Channel struct {
	Name     string
	Ts       int64
	Modes    int
	Clients  map[string]*ChanUser
	Lists    map[int][]string
	Gauge    metrics.Gauge
	Metadata map[string]string

	Lock *sync.Mutex
}

// NewChannel creates a new channel with a given name and ts.
func NewChannel(name string, ts int64) (c *Channel) {
	c = &Channel{
		Name:     name,
		Ts:       ts,
		Lists:    make(map[int][]string),
		Clients:  make(map[string]*ChanUser),
		Modes:    0,
		Gauge:    metrics.NewGauge(),
		Metadata: make(map[string]string),
		Lock:     &sync.Mutex{},
	}

	Etcd.CreateDir("/tetra/channels/"+c.Name, 0)

	Channels[c.Target()] = c

	metrics.Register(strings.ToUpper(name)+"_stats", c.Gauge)

	return
}

// AddChanUser adds a client to the channel, returning the membership.
func (c *Channel) AddChanUser(client *Client) (cu *ChanUser) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	cu = &ChanUser{
		Client:  client,
		Channel: c,
		Prefix:  0,
	}

	c.Clients[client.Uid] = cu

	client.Channels[c.Target()] = c

	return
}

// DelChanUser deletes a client from a channel or returns an error.
func (c *Channel) DelChanUser(client *Client) (err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if _, ok := c.Clients[client.Uid]; !ok {
		err = errors.New("Tried to delete nonexistent chanuser with uid " + client.Uid + " from " + c.Name)
		debug(err)
		return err
	}

	delete(c.Clients, client.Uid)
	delete(client.Channels, c.Name)

	return nil
}

// Target returns a targetable version of Channel.
func (c *Channel) Target() string {
	return strings.ToUpper(c.Name)
}

// IsChannel returns true.
func (c *Channel) IsChannel() bool {
	return true
}
