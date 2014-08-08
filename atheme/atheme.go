// Package atheme implements an Atheme XMLRPC client and does all the
// horrifyingly ugly scraping of the raw output to machine-usable structures.
package atheme

import (
	"github.com/kolo/xmlrpc"
)

// An Atheme context. This contains everything a client needs to access Atheme
// data remotely.
type Atheme struct {
	ServerProxy *xmlrpc.Client
	url         string
	authcookie  string
	Account     string
	ipaddr      string
	Privset     string
}

// Returns a new Atheme instance or raises an error.
func NewAtheme(url string) (atheme *Atheme, err error) {
	var serverproxy *xmlrpc.Client
	serverproxy, err = xmlrpc.NewClient(url, nil)

	if err != nil {
		return nil, err
	}

	atheme = &Atheme{
		ServerProxy: serverproxy,
		url:         url,
		authcookie:  "*",
		Account:     "*",
		ipaddr:      "0",
	}

	return
}
