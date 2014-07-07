package cod

import (
	"fmt"
	"strings"
)

type Client interface {
	Nick() string
	User() string
	Host() string
	Account() string
	Uid() string
	Permissions() int
	Umodes() int
	Ts() int64
	Gecos() string
	Privmsg(source *Client, destination, message string)
	Notice(source *Client, destination, message string)
	Join(name string)
	Euid() string
}

type ServiceClient struct {
	nick        string
	user        string
	host        string
	VHost       string
	Ip          string
	account     string
	uid         string
	gecos       string
	permissions int
	umodes      int
	Commands    map[string]*Command
	Kind        string
	cod         *Cod
	ts          int64
}

func (r ServiceClient) Nick() string {
	return r.nick
}

func (r ServiceClient) User() string {
	return r.user
}

func (r ServiceClient) Host() string {
	return r.VHost
}

func (r ServiceClient) Account() string {
	return r.account
}

func (r ServiceClient) Uid() string {
	return r.uid
}

func (r ServiceClient) Permissions() int {
	return r.permissions
}

func (r ServiceClient) Umodes() int {
	return r.umodes
}

func (r ServiceClient) Gecos() string {
	return r.gecos
}

func (r ServiceClient) Ts() int64 {
	return r.ts
}

func (r ServiceClient) Join(name string) {
	channel := r.cod.Channels[strings.ToLower(name)]

	channel.AddChanUser(Client(r))
}

func (r ServiceClient) Privmsg(source *Client, destination, message string) {}

func (r ServiceClient) Notice(source *Client, destination, message string) {}

func (r ServiceClient) Euid() string {
	return fmt.Sprintf("EUID %s 1 %d +oS %s %s %s %s %s %s :%s", r.nick, r.ts, r.user,
		r.VHost, r.host, r.uid, r.Ip, r.account, r.gecos)
}

type RemoteClient struct {
	nick        string
	user        string
	host        string
	VHost       string
	Ip          string
	account     string
	uid         string
	gecos       string
	permissions int
	umodes      int
	cod         *Cod
	ts          int64
}

func (r RemoteClient) Nick() string {
	return r.nick
}

func (r RemoteClient) User() string {
	return r.user
}

func (r RemoteClient) Host() string {
	return r.VHost
}

func (r RemoteClient) Account() string {
	return r.account
}

func (r RemoteClient) Uid() string {
	return r.uid
}

func (r RemoteClient) Permissions() int {
	return r.permissions
}

func (r RemoteClient) Umodes() int {
	return r.umodes
}

func (r RemoteClient) Gecos() string {
	return r.gecos
}

func (r RemoteClient) message(source Client, kind, destination, message string) {
	r.cod.Conn.SendLine(":%s %s %s :%s", source.Uid(), kind, destination, message)
}

func (r RemoteClient) Privmsg(source *Client, destination, message string) {
	r.message(*source, "PRIVMSG", destination, message)
}

func (r RemoteClient) Notice(source *Client, destination, message string) {
	r.message(*source, "NOTICE", destination, message)
}

func (r RemoteClient) Euid() string {
	return fmt.Sprintf("EUID %s %d + %s %s %s %s %s %s :%s", r.nick, r.user,
		r.VHost, r.host, r.uid, r.Ip, r.gecos)
}

func (r RemoteClient) Join(name string) { }

func (r RemoteClient) Ts() (int64) {
	return r.ts
}

