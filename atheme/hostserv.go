package atheme

import (
	"fmt"
	"strings"
)

// HostServ wraps Atheme's HostServ for use in Go programs.
type HostServ struct {
	a *Atheme
}

// VHost is a vhost listing or request.
type VHost struct {
	Nick  string `json:"nick"`
	VHost string `json:"vhost"`
	Date  string `json:"date"`
}

// String satisfies fmt.Stringer
func (v *VHost) String() string {
	return fmt.Sprintf("nick: %s - vhost: %s - date: %s",
		v.Nick, v.VHost, v.Date)
}

// Activate activates a VHost request for account.
func (hs *HostServ) Activate(account string) (err error) {
	_, err = hs.a.Command("HostServ", "ACTIVATE", account)

	return
}

// Request requests a vhost for the logged in account.
func (hs *HostServ) Request(vhost string) (err error) {
	_, err = hs.a.Command("HostServ", "REQUEST", vhost)

	return
}

// Reject rejects a vhost request for an account.
func (hs *HostServ) Reject(account, message string) (err error) {
	_, err = hs.a.Command("HostServ", "REJECT", account, message)

	return
}

// Revoke revokes a vhost from an account or returns an error.
func (hs *HostServ) Revoke(account string) (err error) {
	_, err = hs.a.Command("HostServ", "VHOST", account)

	return
}

// Assign assigns a vhost to an account or returns an error.
func (hs *HostServ) Assign(account, vhost string) (err error) {
	_, err = hs.a.Command("HostServ", "VHOST", account, vhost)

	return
}

// List returns a list of all the vhosts Atheme is keeping track of.
func (hs *HostServ) List() ([]VHost, error) {
	return hs.ListPattern("*")
}

// ListPattern returns a list of all the vhosts Atheme is keeping track of that
// match a given pattern
func (hs *HostServ) ListPattern(pattern string) (res []VHost, err error) {
	var output string
	output, err = hs.a.Command("HostServ", "LISTVHOST", pattern)
	if err != nil {
		return nil, err
	}

	vhosts := strings.Split(output, "\n")
	vhosts = vhosts[:len(vhosts)-1]

	for _, vhost := range vhosts {
		vhost = strings.Replace(vhost, "  ", "", -1)
		split := strings.Split(vhost, " ")

		res = append(res, VHost{
			Nick:  split[1],
			VHost: split[2],
		})
	}

	return
}

// Waiting returns a list of all the vhosts that are waiting for activation.
func (hs *HostServ) Waiting() (res []VHost, err error) {
	var output string
	output, err = hs.a.Command("HostServ", "WAITING")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// Nick:jewels, vhost:kiss.my.ass.geez (jewels - May 26 17:17:32 2014)
		nick := strings.Split(line, ", ")[0]
		nick = strings.TrimPrefix(nick, "Nick:")
		vhost := strings.Split(line, ", vhost:")[1]
		vhost = strings.Split(vhost, " (")[0]
		date := strings.Split(line, " (")[1]
		date = strings.Split(date, " - ")[1]
		date = strings.TrimSuffix(date, ")")

		res = append(res, VHost{
			Nick:  nick,
			VHost: vhost,
			Date:  date,
		})
	}

	return
}
