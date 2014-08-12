package tetra

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Xe/Tetra/1459"
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type Script struct {
	Name     string
	L        *lua.State
	Tetra    *Tetra
	Log      *log.Logger
	Handlers map[string]*Handler
	Commands map[string]*Command
	Hooks    []*Hook
	Service  string
	Client   *Client
	Uuid     string
	Kind     string
}

func (tetra *Tetra) LoadScript(name string) (script *Script, err error) {
	kind := strings.Split(name, "/")[0]
	client, ok := tetra.Services[kind]
	if !ok {
		client = tetra.Services["tetra"]
	}

	if _, present := tetra.Scripts[name]; present {
		return nil, errors.New("Double script load!")
	}

	script = &Script{
		Name:     name,
		L:        luar.Init(),
		Tetra:    tetra,
		Log:      log.New(os.Stdout, name+" ", log.LstdFlags),
		Handlers: make(map[string]*Handler),
		Commands: make(map[string]*Command),
		Service:  kind,
		Client:   client,
		Uuid:     uuid.New(),
	}

	script.seed()

	script, err = tetra.loadLuaScript(script)
	if err != nil {
		script, err = tetra.loadMoonScript(script)

		if err != nil {
			return nil, errors.New("No such script " + name)
		}
	}

	tetra.Scripts[name] = script

	return
}

func (tetra *Tetra) loadLuaScript(script *Script) (*Script, error) {
	script.L.DoFile("modules/base.lua")

	err := script.L.DoFile("modules/" + script.Name + ".lua")

	if err != nil {
		return script, err
	}

	if tetra.Config.General.Debug || tetra.Bursted {
		script.Log.Printf("lua script %s loaded at %s", script.Name, script.Uuid)
	}

	script.Kind = "lua"

	return script, nil
}

func (tetra *Tetra) loadMoonScript(script *Script) (*Script, error) {
	contents, failed := ioutil.ReadFile("modules/" + script.Name + ".moon")

	if failed != nil {
		return script, errors.New("Could not read " + script.Name + ".moon")
	}

	luar.Register(script.L, "", luar.Map{
		"moonscript_code_from_file": string(contents),
	})

	err := script.L.DoString(`
		moonscript = require "moonscript"

		xpcall = unsafe_xpcall
		pcall = unsafe_pcall

		local func, err = moonscript.loadstring(moonscript_code_from_file)

		if err ~= nil then
			tetra.log.Printf("Moonscript error, %#v", err)
			error(err)
		end

		func()`)
	if err != nil {
		script.Log.Print(err)
		return nil, err
	}

	if tetra.Config.General.Debug || tetra.Bursted {
		script.Log.Printf("moonscript script %s loaded at %s", script.Name, script.Uuid)
	}

	script.Kind = "moonscript"

	return script, nil
}

func (script *Script) seed() {
	luar.Register(script.L, "", luar.Map{
		"client": script.Client,
		"print":  script.Log.Print,
		"script": script,
		"log":    script.Log,
	})

	luar.Register(script.L, "tetra", luar.Map{
		"script":    script,
		"log":       script.Log,
		"bot":       script.Tetra,
		"protohook": script.AddLuaProtohook,
	})

	luar.Register(script.L, "uuid", luar.Map{
		"new": uuid.New,
	})

	luar.Register(script.L, "web", luar.Map{
		"get":  http.Get,
		"post": http.Post,
	})

	luar.Register(script.L, "ioutil", luar.Map{
		"readall":     ioutil.ReadAll,
		"byte2string": byteSliceToString,
	})

	luar.Register(script.L, "strings", luar.Map{
		"join":  strings.Join,
		"split": strings.Split,
		"first": func(str string) string {
			if len(str) > 0 {
				return string(str[0])
			} else {
				return ""
			}
		},
		"rest": func(str string) string {
			if len(str) > 0 {
				return str[1:]
			} else {
				return ""
			}
		},
	})
}

// Call calls a command in a Script.
func (s *Script) Call(command string, source *Client, dest Targeter, args []string) (string, error) {
	cmd, present := s.Client.Commands[command]
	if !present {
		return "", errors.New("No command " + command)
	}

	result := cmd.Impl(source, dest, args)

	return result, nil
}

// AddLuaProtohook adds a lua function as a protocol hook
func (script *Script) AddLuaProtohook(verb string, name string) error {
	function := luar.NewLuaObjectFromName(script.L, name)

	handler, err := script.Tetra.AddHandler(verb, func(line *r1459.RawLine) {
		_, err := function.Call(line)
		if err != nil {
			script.Log.Printf("Lua error %s: %s", script.Name, err.Error())
			script.Log.Printf("%#v", err)
			script.Client.ServicesLog(fmt.Sprintf("%s: %s", script.Name, err.Error()))
		}
	})
	if err != nil {
		return err
	}

	handler.Script = script
	script.Handlers[handler.Uuid] = handler

	return nil
}

// AddLuaCommand adds a new command to a script from a lua context
func (script *Script) AddLuaCommand(verb string, name string) error {
	function := luar.NewLuaObjectFromName(script.L, name)

	command, err := script.Client.NewCommand(verb, func(client *Client, target Targeter, args []string) string {
		reply, err := function.Call(client, target, args)

		if err != nil {
			script.Log.Printf("Lua error %s: %s", script.Name, err.Error())
			script.Log.Printf("%#v", err)
			script.Client.ServicesLog(fmt.Sprintf("%s: %s", script.Name, err.Error()))
			return ""
		}

		return reply.(string)
	})

	if err != nil {
		return err
	}

	command.Script = script

	script.Commands[command.Uuid] = command

	return nil
}

// Adds a hook from lua
func (script *Script) AddLuaHook(verb string, name string) error {
	function := luar.NewLuaObjectFromName(script.L, name)

	hook := script.Tetra.NewHook(verb, func(args ...interface{}) {
		_, err := function.Call(args...)
		if err != nil {
			script.Log.Printf("Lua error %s: %s", script.Name, err.Error())
			script.Log.Printf("%#v", err)
			script.Client.ServicesLog(fmt.Sprintf("%s: %s", script.Name, err.Error()))
		}
	})

	script.Hooks = append(script.Hooks, hook)

	return nil
}

// Unload a script and delete its commands and handlers
func (tetra *Tetra) UnloadScript(name string) error {
	if _, ok := tetra.Scripts[name]; !ok {
		panic("No such script " + name)
	}

	script := tetra.Scripts[name]

	for _, handler := range script.Handlers {
		tetra.DelHandler(handler.Verb, handler.Uuid)
		delete(script.Handlers, handler.Uuid)
	}

	for _, command := range script.Commands {
		delete(script.Commands, command.Uuid)
	}

	for _, hook := range script.Hooks {
		tetra.DelHook(hook)
	}

	script.L.Close()

	delete(tetra.Scripts, name)

	return nil
}

func byteSliceToString(slice []byte) string {
	return string(slice)
}
