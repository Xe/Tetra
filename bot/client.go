package tetra

import (
	"fmt"
	"github.com/Xe/Tetra/bot/modes"
	"strings"
	"time"
)

type Client struct {
	Nick        string `json:"nick"`
	User        string `json:"user"`
	Host        string `json:"realhost"`
	VHost       string `json:"host"`
	Ip          string `json:"ip"`
	Account     string `json:"account"`
	Uid         string `json:"uid"`
	Gecos       string `json:"gecos"`
	Permissions int
	Umodes      int
	Kind        string
	tetra       *Tetra
	Ts          int64               `json:"ts"`
	Channels    map[string]*Channel `json:"channels"`
	Server      *Server             `json:"server"`
}

type Targeter interface {
	Target() string
	IsChannel() bool
}

func (r *Client) Euid() string {
	return fmt.Sprintf("EUID %s 1 %d +oS %s %s %s %s %s %s :%s", r.Nick, r.Ts, r.User,
		r.VHost, r.Host, r.Uid, r.Ip, r.Account, r.Gecos)
}

func (r *Client) Quit() {
	str := fmt.Sprintf(":%s QUIT :Service unloaded")
	r.tetra.Conn.SendLine(str)
}

func (r *Client) message(source *Client, kind string,
	destination Targeter, message string) {
	str := fmt.Sprintf(":%s %s %s :%s", source.Uid, kind, destination.Target(), message)
	r.tetra.Conn.SendLine(str)
}

func (r *Client) Privmsg(destination Targeter, message string) {
	r.message(r, "PRIVMSG", destination, message)
}

func (r *Client) ServicesLog(message string) {
	r.Privmsg(r.tetra.Channels["#SERVICES"], message)
}

func (r *Client) Notice(destination Targeter, message string) {
	r.message(r, "NOTICE", destination, message)
}

func (r *Client) Target() string {
	return r.Uid
}

func (r *Client) IsChannel() bool {
	return false
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

	if r.tetra.Bursted {
		str := fmt.Sprintf(":%s SJOIN %d %s + :%s", r.tetra.Info.Sid, channel.Ts,
			channel.Name, r.Uid)

		r.tetra.Conn.SendLine(str)
	}
}

func (r *Client) Part(channame string) bool {
	upperchan := strings.ToUpper(channame)

	channel, err := r.tetra.Channels[upperchan]
	if !err {
		return err
	}

	channel.DelChanUser(r)

	if r.tetra.Bursted {
		str := fmt.Sprintf(":%s PART %s", r.Uid, channame)

		r.tetra.Conn.SendLine(str)
	}

	return true
}

func (r *Client) IsOper() bool {
	return r.Umodes&modes.UPROP_IRCOP == modes.UPROP_IRCOP
}
