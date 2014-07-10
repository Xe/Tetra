package tetra

import (
	"errors"
	"fmt"
	"strconv"
)

type Client struct {
	Nick        string
	User        string
	Host        string
	VHost       string
	Ip          string
	Account     string
	Uid         string
	Gecos       string
	Permissions int
	Umodes      int
	Commands    map[string]*Command
	Kind        string
	tetra       *Tetra
	Ts          int64
}

type Targeter interface {
	Target() string
}

func (r *Client) Euid() string {
	return fmt.Sprintf("EUID %s 1 %d +oS %s %s %s %s %s %s :%s", r.Nick, r.Ts, r.User,
		r.VHost, r.Host, r.Uid, r.Ip, r.Account, r.Gecos)
}

func (r *Client) message(source *Client, kind string,  destination Targeter, message string) {
	r.tetra.Conn.SendLine(":%s %s %s :%s", source.Uid, kind, destination, message)
}

func (r *Client) Privmsg(destination Targeter, message string) {
	r.message(r, "PRIVMSG", destination, message)
}

func (r *Client) Notice(destination Targeter, message string) {
	r.message(r, "NOTICE", destination, message)
}

func (r *Client) Target() (string) {
	return r.Uid
}

func (r *Client) Join(channel *Channel) (err error) {
	if r == nil {
		panic(errors.New("What the fuck"))
	}

	r.tetra.Conn.SendLine(":420 SJOIN " + strconv.FormatInt(channel.Ts, 10) + " " +
		channel.Name + " + :@" + r.Uid)

	return
}

