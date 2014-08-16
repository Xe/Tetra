package atheme

import (
	"strings"
)

const (
	NICKSERV_HOLD       = "Hold"       // Account cannot expire
	NICKSERV_HIDEMAIL   = "HideMail"   // Account email hidden
	NICKSERV_NEVEROP    = "NeverOp"    // Account can't be added to access lists
	NICKSERV_NOOP       = "NoOp"       // Account can't be opped by services
	NICKSERV_NOMEMO     = "NoMemo"     // Account cannot receive memos
	NICKSERV_EMAILMEMOS = "EMailMemos" // Account has memos emailed to it
	NICKSERV_PRIVATE    = "Private"    // Account information is private
)

// Struct NickServ implements a Golang client to Atheme's NickServ. This is
// mostly a port of Cod's Atheme parsing code
type NickServ struct {
	a Atheme
}

type Flagset []map[string]string

func (ns *NickServ) parseAccess(data string) (res Flagset) {
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

// ListOwnAccess lists the channels a user has flags in.
func (ns *NickServ) ListOwnAccess() (res Flagset, err error) {
	var temp string
	temp, err = ns.a.Command("NickServ", "LISTCHANS")
	if err != nil {
		return nil, err
	}

	res = ns.parseAccess(temp)

	return
}
