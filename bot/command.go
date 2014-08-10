package tetra

import (
	"errors"
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

// Struct command holds everything needed for a bot command.
type Command struct {
	Impl      func(*Client, Targeter, []string) string
	Uuid      string
	Script    *Script
	Verb      string
	Owner     *Client
	NeedsOper bool
}

// NewCommand returns a new command instance.
func (c *Client) NewCommand(verb string, handler func(*Client, Targeter, []string) string) (cmd *Command, err error) {
	verb = strings.ToUpper(verb)

	if _, present := c.Commands[verb]; present {
		return nil, errors.New("Duplicate command " + verb + " for client " + c.Nick)
	}

	cmd = &Command{
		Impl:  handler,
		Owner: c,
		Verb:  verb,
		Uuid:  uuid.New(),
	}

	c.Commands[verb] = cmd

	return
}
