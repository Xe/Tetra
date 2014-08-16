// Package atheme implements an Atheme XMLRPC client and does all the
// horrifyingly ugly scraping of the raw output to machine-usable structures.
package atheme

import (
	"strings"

	"github.com/Xe/Tetra/atheme/xmlrpc"
)

// An Atheme context. This contains everything a client needs to access Atheme
// data remotely.
type Atheme struct {
	Privset     []string // Privilege set of the user
	Account     string   // Account Atheme is logged in as
	serverProxy *xmlrpc.Client
	url         string
	Authcookie  string
	ipaddr      string
	NickServ    *NickServ
	ChanServ    *ChanServ
}

// Returns a new Atheme instance or raises an error.
func NewAtheme(url string) (atheme *Atheme, err error) {
	var serverproxy *xmlrpc.Client
	serverproxy, err = xmlrpc.NewClient(url, nil)

	if err != nil {
		return nil, err
	}

	atheme = &Atheme{
		Account:     "*",
		serverProxy: serverproxy,
		url:         url,
		Authcookie:  "*",
		ipaddr:      "0",
		NickServ:    &NickServ{a: atheme},
		ChanServ:    &ChanServ{a: atheme},
	}

	return
}

// Command runs an Atheme command and gives the output or an error.
func (a *Atheme) Command(args ...string) (res string, err error) {
	err = a.serverProxy.Call("atheme.command", args, &res)

	return
}

// Login attempts to log a user into Atheme. It returns true or false
func (a *Atheme) Login(username, password string) (success bool, err error) {
	var authcookie string

	err = a.serverProxy.Call("atheme.login", []string{username, password, "::1"}, &authcookie)

	if err == nil {
		a.Authcookie = authcookie
		a.Account = username
		success = true
	} else {
		return false, err
	}

	return
}

// Logout logs a user out of Atheme. There is no return.
func (a *Atheme) Logout() {
	var res string

	a.serverProxy.Call("atheme.logout", []string{a.Authcookie, a.Account}, &res)

	a.Account = "*"
	a.Authcookie = "*"

	return
}

// GetPrivset returns the privset of a user.
func (a *Atheme) GetPrivset() (privs []string) {
	if a.Privset == nil {
		var res string

		a.serverProxy.Call("atheme.privset", []string{a.Authcookie, a.Account}, &res)

		a.Privset = strings.Split(res, " ")
	}

	return a.Privset
}
