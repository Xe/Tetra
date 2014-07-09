package atheme

import (
	"github.com/kolo/xmlrpc"
)

type Atheme struct {
	ServerProxy *xmlrpc.Client
	url         string
	authcookie  string
	Account     string
	ipaddr      string
	Privset     string
}

func NewAtheme(url string) (atheme *Atheme, err error) {
	var serverproxy *xmlrpc.Client
	serverproxy, err = xmlrpc.NewClient(url, nil)

	if err != nil{
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
