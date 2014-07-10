package tetra

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/Xe/Tetra/1459"
)

type Handler struct {
	Impl   func(*r1459.RawLine)
	Verb   string
	Uuid   string
	Script *Script
}

func (tetra *Tetra) AddHandler(verb string, impl func(*r1459.RawLine)) (handler *Handler, err error) {
	handler = &Handler{
		Verb: verb,
		Impl: impl,
		Uuid: uuid.New(),
	}

	if _, ok := tetra.Handlers[verb]; !ok {
		tetra.Handlers[verb] = make(map[string]*Handler)
	}

	tetra.Handlers[verb][handler.Uuid] = handler

	return
}

func (tetra *Tetra) DelHandler(verb string, uuid string) (err error) {
	if _, present := tetra.Handlers[verb]; !present {
		return errors.New("No such verb to delete handler for " + verb)
	}

	delete(tetra.Handlers[verb], uuid)

	return nil
}
