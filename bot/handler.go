package tetra

import (
	"errors"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Xe/Tetra/1459"
)

// Struct Handler defines a raw protocol verb handler. Please do not use this
// unless you have good reason to.
type Handler struct {
	Impl   func(*r1459.RawLine)
	Verb   string
	Uuid   string
	Script *Script
	Go     bool
}

// AddHandler adds a handler for a given verb.
func AddHandler(verb string, impl func(*r1459.RawLine)) (handler *Handler, err error) {
	handler = &Handler{
		Verb: verb,
		Impl: impl,
		Uuid: uuid.New(),
		Go:   true,
	}

	if _, ok := Handlers[verb]; !ok {
		Handlers[verb] = make(map[string]*Handler)
	}

	Handlers[verb][handler.Uuid] = handler

	return
}

// DelHandler deletes a handler for a given protocol verb by the UUID of the handler.
func DelHandler(verb string, uuid string) (err error) {
	if _, present := Handlers[verb]; !present {
		err = errors.New("No such verb to delete handler for " + verb)
		debug(err)
		return err
	}

	Handlers[verb][uuid].Go = false

	delete(Handlers[verb], uuid)

	return nil
}
