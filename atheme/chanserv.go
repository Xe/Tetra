package atheme

import (
	"strconv"
	"strings"
)

// Struct ChanServ implements a Golang client to Atheme's ChanServ. This is
// mostly a port of Cod's string parsing code.
type ChanServ struct {
	a *Atheme
}

// Struct Flagset is a simple flagset wrapper.
type Flagset struct {
	Id    int    `json:"id"`
	Nick  string `json:"nick"`
	Flags string `json:"flags"`
}

// Struct ChannelInfo is the information Atheme has on a channel.
type ChannelInfo struct {
	Name   string   `json:"name"`   // Channel name
	Mlock  string   `json:"mlock"`  // Channel mode lock
	Flags  []string `json:"flags"`  // Channel SET flags
	Prefix string   `json:"prefix"` // Channel FANTASY prefix
}

// Kick sends a ChanServ KICK command to channel on victim with the denoted
// reason. You must have a reason for calls made with this function.
func (cs *ChanServ) Kick(channel, victim, reason string) (res string, err error) {
	return cs.a.Command("ChanServ", "KICK", channel, victim, reason)
}

// GetAccessList returns a slice of Flagsets representing the access
// list of the channel you are requesting. This will fail if the Atheme call
// fails.
func (cs *ChanServ) GetAccessList(channel string) (res []Flagset, err error) {
	var output string

	output, err = cs.a.Command("ChanServ", "FLAGS", channel)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.Replace(line, "  ", "", -1)
		data := strings.Split(line, " ")

		id, _ := strconv.Atoi(data[0])

		res = append(res, Flagset{
			Id:    id,
			Nick:  data[1],
			Flags: data[2],
		})
	}

	return
}

// SetAccessList commits a flag change on a channel with a given flagset.
func (cs *ChanServ) SetAccessList(channel, target, flags string) (err error) {
	_, err = cs.a.Command("ChanServ", "FLAGS", channel, target, flags)

	return
}

// GetChannelInfo gets information on a channel, returning a ChannelInfo struct
// or an error describing the fault.
func (cs *ChanServ) GetChannelInfo(channel string) (ci *ChannelInfo, err error) {
	// I am sorry.
	var output string
	output, err = cs.a.Command("ChanServ", "INFO", channel, "FOO")
	if err != nil {
		return nil, err
	}

	/*
		Information on #niichan:
		Registered : Nov 06 10:01:32 2013 (40w 2d 12h ago)
		Mode lock  : +n
		Flags      : HOLD SECURE VERBOSE KEEPTOPIC GUARD FANTASY PRIVATE
		Prefix     : ! (default)
		*** End of Info ***
	*/

	ci = &ChannelInfo{
		Name: channel,
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "Information on #") || strings.HasPrefix(line, "*") {
			continue
		}

		line = strings.Replace(line, "  ", "", -1)
		data := strings.Split(line, ":")
		key := strings.ToLower(data[0])
		value := strings.TrimSpace(data[1])

		// TODO: replace this with reflect
		switch key {
		case "mode lock":
			ci.Mlock = value
		case "flags":
			ci.Flags = strings.Split(value, " ")
		case "prefix":
			ci.Prefix = value
		}
	}

	return
}

// GetChannelFlags returns the SET flags of a channel as a string slice.
func (cs *ChanServ) GetChannelFlags(channel string) (flags []string, err error) {
	var ci *ChannelInfo

	ci, err = cs.GetChannelInfo(channel)
	if err != nil {
		return nil, err
	}

	return ci.Flags, err
}

// SetChannelFlag sets a channel SET flag or returns an error.
func (cs *ChanServ) SetChannelFlag(channel, flag, value string) (err error) {
	_, err = cs.a.Command("ChanServ", "SET", channel, flag, value)

	return
}
