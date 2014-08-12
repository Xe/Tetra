package tetra

import (
	"fmt"
	"github.com/Xe/Tetra/bot/modes"
	"strings"
	"time"
)

// Struct Client holds information about a client on the IRC network.
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
	Channels    map[string]*Channel
	Server      *Server
	Commands    map[string]*Command
	Certfp      string
}

// Interface Targeter wraps around Client and Channel to make messaging to them
// seamless.
type Targeter interface {
	Target() string  // Targetable version of name
	IsChannel() bool // Is this a channel?
}

// Euid returns an EUID burst.
func (r *Client) Euid() string {
	return fmt.Sprintf("EUID %s 1 %d +oS %s %s %s %s %s %s :%s", r.Nick, r.Ts, r.User,
		r.VHost, r.Host, r.Uid, r.Ip, r.Account, r.Gecos)
}

// Quit quits a client off of the network.
func (r *Client) Quit() {
	str := fmt.Sprintf(":%s QUIT :Service unloaded", r.Uid)
	r.tetra.Conn.SendLine(str)
}

func (r *Client) message(source *Client, kind string, destination Targeter, message string) {
	if message == "" {
		message = " "
	}

	str := fmt.Sprintf(":%s %s %s :%s", source.Uid, kind, destination.Target(), message)
	r.tetra.Conn.SendLine(str)
}

// Privmsg sends a PRIVMSG to destination with given message.
func (r *Client) Privmsg(destination Targeter, message string) {
	r.message(r, "PRIVMSG", destination, message)
}

// ServicesLog logs a given message to the services snoop channel.
func (r *Client) ServicesLog(message string) {
	r.Privmsg(r.tetra.Channels[strings.ToUpper(r.tetra.Config.General.SnoopChan)], message)
}

// Notice sends a NOTICE to destination with given message.
func (r *Client) Notice(destination Targeter, message string) {
	r.message(r, "NOTICE", destination, message)
}

// Target returns a targetable version of a Client.
func (r *Client) Target() string {
	return r.Uid
}

// IsChannel returns false.
func (r *Client) IsChannel() bool {
	return false
}

// Join makes the client join a channel. This does not check bans.
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

// Part makes the client leave a channel.
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

// IsOper returns if the client is an operator or not.
func (r *Client) IsOper() bool {
	return r.Umodes&modes.UPROP_IRCOP == modes.UPROP_IRCOP
}
