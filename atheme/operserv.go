package atheme

import (
	"fmt"
	"strconv"
	"strings"
)

// OperServ models Atheme's OperServ to Golang programs.
type OperServ struct {
	a *Atheme
}

// Akill models an AKILL, also known as a K:Line. This represents a network-wide
// ban by user@host.
type Akill struct {
	Num    int    `json:"num"`
	Mask   string `json:"mask"`
	Setter string `json:"setter"`
	Expiry string `json:"expiry"`
	Reason string `json:"reason"`
}

// AkillAdd adds an Akill to Atheme and returns the struct created by the call.
func (os *OperServ) AkillAdd(mask, reason, time string) (ak *Akill, err error) {
	_, err = os.a.Command("OperServ", "AKILL", "ADD", mask, "!T", time, reason)
	if err != nil {
		return nil, err
	}

	ak = &Akill{
		Mask:   mask,
		Reason: reason,
		Setter: os.a.Account,
	}

	return
}

// AkillDel attempts to delete an Akill from Atheme and returns an error
// representing the status of the command.
func (os *OperServ) AkillDel(num int) (err error) {
	_, err = os.a.Command("OperServ", "AKILL", "DEL", fmt.Sprintf("%d", num))

	return
}

// AkillList returns the Akills Atheme is tracking and an error representing
// the command's status.
func (os *OperServ) AkillList() (akills []Akill, err error) {
	var output string
	output, err = os.a.Command("OperServ", "AKILL", "LIST", "FULL")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if line[0] == 'A' || line[0] == 'T' {
			continue
		}

		// 72: *@8.8.8.8 - by Xe - expires in 0 days, 0:22:07 - (test)
		numbersplit := strings.Split(line, ": ")
		num, _ := strconv.Atoi(numbersplit[0])

		// *@8.8.8.8 - by Xe - expires in 0 days, 0:22:07 - (test)
		infosplit := strings.Split(numbersplit[1], " - ")
		mask := infosplit[0]
		setter := strings.Split(infosplit[1], " ")[1]
		expiry := strings.TrimPrefix(infosplit[2], "expires in ")

		reason := infosplit[3]
		reason = strings.TrimPrefix(reason, "(")
		reason = strings.TrimSuffix(reason, ")")

		akills = append(akills, Akill{
			Num:    num,
			Setter: setter,
			Mask:   mask,
			Reason: reason,
			Expiry: expiry,
		})
	}

	return
}

// Kill asks OperServ to KILL a client off of the network for a given reason.
// You must give a reason.
func (os *OperServ) Kill(target, reason string) (err error) {
	_, err = os.a.Command("OperServ", "KILL", target, reason)

	return
}
