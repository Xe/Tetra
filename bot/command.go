package cod

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
)

type Command struct {
	Impl   func(*RemoteClient, []string)
	Uuid   string
	Script Script
	Owner  *ServiceClient
	Verb   string
	Perms  int
}

func (cod *Cod) AddCommand(service, verb string, impl func(*RemoteClient, []string)) (command *Command, err error) {
	client := cod.Services[service]

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

func (cod *Cod) DelCommand(service, verb string) (err error) {
	client := cod.Services[service]

	if _, present := client.Commands[verb]; present {
		return errors.New("No such command " + verb + " for service " + service)
	}

	delete(client.Commands, verb)

	return
}
