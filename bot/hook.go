package tetra

import (
	"errors"
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

// Struct Hook defines a command hook for Tetra. This can be used for hooking on
// events (like being yo'd).
type Hook struct {
	Uuid  string
	impl  func(...interface{})
	Owner *Script
	Verb  string
}

// NewHook allocates and returns a new Hook structure
func (t *Tetra) NewHook(verb string, impl func(...interface{})) (h *Hook) {
	verb = strings.ToUpper(verb)

	h = &Hook {
		Uuid: uuid.New(),
		Verb: verb,
		impl: impl,
	}

	t.Hooks[verb] = append(t.Hooks[verb], h)

	return
}

// RunHook runs a hook in parallel across multiple goroutines, one per implementaion
// of the hook
func (t *Tetra) RunHook(verb string, args ...interface{}) (err error) {
	if _, present := t.Hooks[verb]; present {
		for _, hook := range t.Hooks[verb] {
			go hook.impl(args...)
		}
	} else {
		return errors.New("No such hook " + verb)
	}

	return
}

// DelHook deletes a hook.
func (t *Tetra) DelHook(hook *Hook) (err error) {
	if _, present := t.Hooks[hook.Verb]; !present {
		return errors.New("Improper hook.")
	}

	var target int

	for _, myhook := range t.Hooks[hook.Verb] {
		if hook.Uuid == myhook.Uuid {
			break
		}

		target++
	}

	t.Hooks[hook.Verb] = append(t.Hooks[hook.Verb][:target], t.Hooks[hook.Verb][target+1:]...)

	return
}
