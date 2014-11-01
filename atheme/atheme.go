// Package atheme implements an Atheme XMLRPC client and does all the
// horrifyingly ugly scraping of the raw output to machine-usable structures.
package atheme

import (
	"strings"
	"time"

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
	OperServ    *OperServ
	HostServ    *HostServ
	MemoServ    *MemoServ
	LastUsed    time.Time // When the last RPC call was made
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
		LastUsed:    time.Now(),
	}

	atheme.NickServ = &NickServ{a: atheme}
	atheme.ChanServ = &ChanServ{a: atheme}
	atheme.OperServ = &OperServ{a: atheme}
	atheme.HostServ = &HostServ{a: atheme}
	atheme.MemoServ = &MemoServ{a: atheme}

	return atheme, nil
}

// Command runs an Atheme command and gives the output or an error.
func (a *Atheme) Command(args ...string) (string, error) {
	var result string

	fullcommand := []string{a.Authcookie, a.Account, a.ipaddr}

	for _, arg := range args {
		fullcommand = append(fullcommand, arg)
	}

	err := a.serverProxy.Call("atheme.command", &fullcommand, &result)

	a.LastUsed = time.Now()

	return result, err
}

// Login attempts to log a user into Atheme. It returns true or false
func (a *Atheme) Login(username, password string) (err error) {
	var authcookie string

	err = a.serverProxy.Call("atheme.login", []string{username, password, "::1"}, &authcookie)

	if err != nil {
		return err
	}

	a.Authcookie = authcookie
	a.Account = username

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
