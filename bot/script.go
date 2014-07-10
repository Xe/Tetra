package tetra

import (
	"code.google.com/p/go-uuid/uuid"
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"github.com/Xe/Tetra/1459"
	"log"
	"os"
	"strings"
)

type Script struct {
	Name     string
	L        *lua.State
	Tetra    *Tetra
	Log      *log.Logger
	Commands []*Command
	Handlers []*Handler
	Service  string
}

func (tetra *Tetra) LoadScript(name string) (script *Script) {
	script = &Script{
		Name:     name,
		L:        luar.Init(),
		Tetra:    tetra,
		Log:      log.New(os.Stdout, name+" ", log.LstdFlags),
		Commands: nil,
		Handlers: nil,
		Service:  strings.Split(name, "/")[0],
	}

	luar.Register(script.L, "", luar.Map{
		"service": script.Service,
	})

	luar.Register(script.L, "tetra", luar.Map{
		"script": script,
		"log":    script.Log,
		"bot":    tetra,
	})

	luar.Register(script.L, "uuid", luar.Map{
		"new": uuid.New,
	})

	script.L.DoFile("modules/base.lua")

	tetra.Scripts[name] = script

	err := script.L.DoFile("modules/" + name + ".lua")
	if err != nil {
		script.Log.Printf("%#v", err)
	}

	return
}

// Add a lua command (by name) from a lua script. This is designed to be ran
// from a lua environment.
func (script *Script) AddLuaCommand(tetra *Tetra, verb string, help string, funcname string) {
	function := luar.NewLuaObjectFromName(script.L, funcname)

	command, _ := tetra.AddCommand(script.Service, verb,
		func(client *Client, message []string) string {
			reply, err := function.Call(client, message)
			if err != nil {
				script.Log.Printf("Lua error %s: %#v", script.Name, err)
				return "Lua error"
			}

			return reply.(string)
		})

	script.Commands = append(script.Commands, command)
}

// Add a lua function as a protocol hook
func (script *Script) AddLuaProtohook(tetra *Tetra, verb string, name string) {
	function := luar.NewLuaObjectFromName(script.L, name)

	handler, err := tetra.AddHandler(verb, func(line *r1459.RawLine) {
		_, err := function.Call(line)
		if err != nil {
			script.Log.Printf("Lua error %s: %#v", script.Name, err)
		}
	})
	if err != nil {
		panic(err)
	}

	handler.Script = script
	script.Handlers = append(script.Handlers, handler)
}

// Unload a script and delete its commands and handlers
func (tetra *Tetra) UnloadScript(name string) error {
	if _, ok := tetra.Scripts[name]; !ok {
		panic("No such script " + name)
	}

	script := tetra.Scripts[name]

	for i, command := range script.Commands {
		tetra.DelCommand(command.Verb, command.Uuid)
		script.Commands = script.Commands[i:]
	}

	for j, handler := range script.Handlers {
		tetra.DelHandler(handler.Verb, handler.Uuid)
		script.Handlers = script.Handlers[j:]
	}

	script.L.Close()

	delete(tetra.Scripts, name)

	return nil
}
