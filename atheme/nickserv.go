package atheme

import (
	"fmt"
	"strings"
)

// Struct NickServ implements a Golang client to Atheme's NickServ. This is
// mostly a port of Cod's Atheme parsing code
type NickServ struct {
	a *Atheme
}

// NickServFlagset is a convenience wrapper around a slice of string->string
// maps.
type NickServFlagset []map[string]string

// parseAccess parses the access strings from Atheme and returns a nice usable
// NickServFlagset.
func (ns *NickServ) parseAccess(data string) (res NickServFlagset) {
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		fields := strings.Split(line, " ")
		appd := make(map[string]string)

		if fields[0] != "Access" {
			continue
		}

		if len(fields) < 4 {
			continue
		}

		appd["channel"] = fields[4]
		appd["flags"] = fields[2]
		res = append(res, appd)
	}

	return
}

// OwnInfo gets NickServ info for the "local" user in the Atheme instance.
func (ns *NickServ) OwnInfo() (map[string]string, error) {
	return ns.Info(ns.a.Account)
}

// Info gets NickServ info on an arbitrary user or returns an error.
func (ns *NickServ) Info(target string) (res map[string]string, err error) {
	var output string
	output, err = ns.a.Command("NickServ", "INFO", target)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Information on") {
			//Information on fooicus (account fooicus):
			split := strings.Split(line, " ")[5]
			accname := split[:len(split)-2]
			res["account"] = accname

			continue
		}

		/*
			Registered : Aug 15 19:44:03 2014 (2h 43m 55s ago)
			Entity ID  : AAAAAAAEC
			Last addr  : xena@62-059-087-073.tukw.qwest.net
			Last seen  : Aug 15 20:34:53 2014 (1h 53m 5s ago)
			User seen  : Aug 15 22:14:19 2014 (13m 39s ago)
			Flags      : HideMail
			Last quit  : Quit: Xaric: If you have a better quit message then submit a patch!
		*/
		fields := strings.Split(line, ":")
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		key = strings.Join(strings.Split(key, " "), "-")

		if key == "Metadata" {
			/*
				Metadata   : friendcode = 0877-1711-6824
				Metadata   : PGP = 0xF637E333
			*/
			metadata := strings.Split(value, " = ")
			key = "meta-" + metadata[0]
			value = metadata[1]
		}

		res[strings.ToLower(key)] = strings.TrimSpace(value)
	}

	return
}

// ListOwnAccess lists the channels this user has flags in.
func (ns *NickServ) ListOwnAccess() (res NickServFlagset, err error) {
	var temp string
	temp, err = ns.a.Command("NickServ", "LISTCHANS")
	if err != nil {
		return nil, err
	}

	res = ns.parseAccess(temp)

	return
}

// ListAccess lists the channels a user has flags in.
func (ns *NickServ) ListAccess(target string) (res NickServFlagset, err error) {
	var temp string
	temp, err = ns.a.Command("NickServ", "LISTCHANS", target)
	if err != nil {
		return nil, err
	}

	res = ns.parseAccess(temp)

	return
}

// Uid returns the Atheme UID of an account.
func (ns *NickServ) Uid(account string) (res string, err error) {
	uidstring, err := ns.a.Command("NickServ", "ACC", account)
	if err != nil {
		return "", err
	}

	if len(strings.Split(uidstring, " ")) != 4 {
		return "", fmt.Errorf("Insufficient permissions")
	}

	return strings.Split(uidstring, " ")[3], nil
}

// SetPassword sets the password for an account.
func (ns *NickServ) SetPassword(password string) (res string, err error) {
	res, err = ns.a.Command("NickServ", "SET", "PASSWORD", password)

	return
}

// SetEmail sets the email address for an account.
func (ns *NickServ) SetEmail(email string) (res string, err error) {
	res, err = ns.a.Command("NickServ", "SET", "EMAIL", email)

	return
}
