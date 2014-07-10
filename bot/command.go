package tetra

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
)

type Command struct {
	Impl   func(*Client, []string) string
	Uuid   string
	Script Script
	Owner  *Client
	Verb   string
	Perms  int
}

func (tetra *Tetra) AddCommand(service, verb string, impl func(*Client, []string) string)  (command *Command, err error) {
	client := tetra.Services[service]

	command = &Command{
		Impl:  impl,
		Owner: client,
		Perms: 0,
		Verb:  verb,
		Uuid:  uuid.New(),
	}

	if _, present := client.Commands[verb]; present {
		err = errors.New("Double command! " + fmt.Sprintf("%#v", command))
		return nil, err
	}

	client.Commands[verb] = command

	return
}

func (tetra *Tetra) DelCommand(service, verb string) (err error) {
	client := tetra.Services[service]

	if _, present := client.Commands[verb]; present {
		return errors.New("No such command " + verb + " for service " + service)
	}

	delete(client.Commands, verb)

	return
}
