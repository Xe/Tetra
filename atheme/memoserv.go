package atheme

import (
	"strconv"
	"strings"
)

// MemoServ binds MemoServ RPC calls.
type MemoServ struct {
	a *Atheme
}

// Memo represents a MemoServ memo.
type Memo struct {
	From    string `json:"from"`
	Date    string `json:"date"`
	Message string `json:"message"`
	ID      int    `json:"id"`
}

// Send figures out what send command to use and sends a memo to that target.
// It returns an error if the command failed.
func (ms *MemoServ) Send(target, message string) (err error) {
	switch target[0] {
	case '!':
		_, err = ms.a.Command("MemoServ", "SENDGROUP", target, message)
	case '#':
		_, err = ms.a.Command("MemoServ", "SENDOPS", target, message)
	default:
		_, err = ms.a.Command("MemoServ", "SEND", target, message)
	}
	return
}

// List lists all the Memos in a user's inbox.
func (ms *MemoServ) List() (memos []Memo, err error) {
	var output string

	lines := strings.Split(output, "\n")

	/*
		You have 3 memos (0 new).

		- 1 From: Quora Sent: Apr 02 15:04:48 2014
		- 2 From: Quora Sent: Apr 05 08:45:15 2014
		- 3 From: Xe Sent: Apr 06 09:12:05 2014
	*/

	if len(lines) < 2 {
		return
	}

	for _, line := range lines[2:] {
		split := strings.Split(line, " ")
		from := split[3]
		id, _ := strconv.Atoi(split[1])
		date := strings.Join(split[5:], " ")

		memos = append(memos, Memo{
			From: from,
			Date: date,
			ID:   id,
		})
	}

	return
}

// Forward forwards a Memo to another account.
func (ms *MemoServ) Forward(memo Memo, target string) (err error) {
	_, err = ms.a.Command("MemoServ", "FORWARD", target, strconv.Itoa(memo.ID))

	return
}
