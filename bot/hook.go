package tetra

import (
	"errors"
	"log"
	"strings"
	"sync"

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
func NewHook(verb string, impl func(...interface{})) (h *Hook) {
	verb = strings.ToUpper(verb)

	h = &Hook{
		Uuid: uuid.New(),
		Verb: verb,
		impl: impl,
	}

	Hooks[verb] = append(Hooks[verb], h)

	return
}

// RunHook runs a hook in parallel across multiple goroutines, one per implementaion
// of the hook. Returns error if there is no such hook.
func RunHook(verb string, args ...interface{}) (err error) {
	debugf("Running hooks for %s", verb)

	if _, present := Hooks[verb]; present {
		wg := sync.WaitGroup{}

		wg.Add(len(Hooks[verb]))

		for i, arg := range args {
			if arg == nil {
				log.Printf(
					"Wtf, arg %d is nil for hook %s args %#v",
					i,
					verb,
					args,
				)
			}
		}

		for _, hook := range Hooks[verb] {
			hook := hook
			go func() {
				hook.impl(args...)
				wg.Done()
			}()
		}

		wg.Wait()
	} else {
		return errors.New("No such hook " + verb)
	}

	return
}

// DelHook deletes a hook. Returns an error if there is no such hook to delete.
func DelHook(hook *Hook) (err error) {
	if _, present := Hooks[hook.Verb]; !present {
		return errors.New("Improper hook.")
	}

	var target int

	for _, myhook := range Hooks[hook.Verb] {
		if hook.Uuid == myhook.Uuid {
			break
		}

		target++
	}

	Hooks[hook.Verb] = append(Hooks[hook.Verb][:target], Hooks[hook.Verb][target+1:]...)

	return
}
