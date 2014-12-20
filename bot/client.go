package tetra

import (
	"fmt"
	"strings"
	"time"

	"github.com/Xe/Tetra/bot/modes"
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
	Ts          int64
	Channels    map[string]*Channel
	Server      *Server
	Commands    map[string]*Command
	Certfp      string
	Metadata    map[string]string
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
	Conn.SendLine(str)
}

func (r *Client) message(source *Client, kind string, destination Targeter, message string) {
	if message == "" {
		message = " "
	}

	str := fmt.Sprintf(":%s %s %s :%s", source.Uid, kind, destination.Target(), message)
	Conn.SendLine(str)
}

// Privmsg sends a PRIVMSG to destination with given message.
func (r *Client) Privmsg(destination Targeter, message string) {
	r.message(r, "PRIVMSG", destination, message)
}

// ServicesLog logs a given message to the services snoop channel.
func (r *Client) ServicesLog(message string) {
	Log.Printf("%s: %s", r.Nick, message)
	r.Privmsg(Channels[strings.ToUpper(ActiveConfig.General.SnoopChan)], message)
}

// OperLog logs a given message to the operator channel.
func (r *Client) OperLog(message string) {
	Log.Printf("%s: %s", r.Nick, message)
	r.Privmsg(Channels[strings.ToUpper(ActiveConfig.General.StaffChan)], message)
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

	if _, ok := Channels[upperchan]; !ok {
		channel = NewChannel(channame, time.Now().Unix())
	} else {
		channel = Channels[upperchan]
	}

	channel.AddChanUser(r)

	if Bursted {
		str := fmt.Sprintf(":%s SJOIN %d %s + :%s", Info.Sid, channel.Ts,
			channel.Name, r.Uid)

		Conn.SendLine(str)
	}
}

// Part makes the client leave a channel.
func (r *Client) Part(channame string) bool {
	upperchan := strings.ToUpper(channame)

	channel, err := Channels[upperchan]
	if !err {
		return err
	}

	channel.DelChanUser(r)

	if Bursted {
		str := fmt.Sprintf(":%s PART %s", r.Uid, channame)

		Conn.SendLine(str)
	}

	return true
}

// IsOper returns if the client is an operator or not.
func (r *Client) IsOper() bool {
	return r.Umodes&modes.UPROP_IRCOP == modes.UPROP_IRCOP
}

// Kill kills a target client
func (r *Client) Kill(target *Client, reason string) {
	str := fmt.Sprintf(":%s KILL %s :%s!%s (%s)", r.Uid, target.Uid,
		ActiveConfig.Server.Name, r.Nick, reason)

	Conn.SendLine(str)
	Clients.DelClient(target)
}

// Chghost changes a client's visible host
func (r *Client) Chghost(target *Client, newhost string) (err error) {
	strings.Replace(newhost, "_", "--", -1)

	target.VHost = newhost

	line := fmt.Sprintf(":%s CHGHOST %s %s", r.Server.Sid, target.Target(), newhost)

	Conn.SendLine(line)

	return
}
