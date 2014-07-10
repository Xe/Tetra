package tetra

import (
	"fmt"
	"github.com/Xe/Tetra/modes"
	"strings"
	"time"
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

func (r *Client) message(source *Client, kind string,
	destination Targeter, message string) {
	str := fmt.Sprintf(":%s %s %s :%s", source.Uid, kind, destination.Target(), message)
	r.tetra.Conn.SendLine(str)
}

func (r *Client) Privmsg(destination Targeter, message string) {
	r.message(r, "PRIVMSG", destination, message)
}

func (r *Client) Notice(destination Targeter, message string) {
	r.message(r, "NOTICE", destination, message)
}

func (r *Client) Target() string {
	return r.Uid
}

func (r *Client) Join(channame string) {
	var channel *Channel

	upperchan := strings.ToUpper(channame)

	if _, ok := r.tetra.Channels[upperchan]; !ok {
		channel = r.tetra.NewChannel(channame, time.Now().Unix())
	} else {
		channel = r.tetra.Channels[upperchan]
	}

	channel.AddChanUser(r)

	str := fmt.Sprintf(":%s SJOIN %d %s + :%s", r.tetra.Info.Sid, channel.Ts,
		channel.Name, r.Uid)

	r.tetra.Conn.SendLine(str)
}

func (r *Client) IsOper() bool {
	return r.Umodes & modes.UPROP_IRCOP == modes.UPROP_IRCOP
}
