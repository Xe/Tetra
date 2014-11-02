package tetra

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Xe/Tetra/1459"
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

// Struct Script implements a Lua scripting interface to Tetra.
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
	Trigger  chan []interface{}
}

// The different kinds of invocations that can be called upon.
const (
	INV_COMMAND = 0x0001
	INV_NAMHOOK = 0x0002
	INV_PROHOOK = 0x0004
)

// Struct Invocation represents an event from Go->Lua.
type Invocation struct {
	Kind     int
	Args     []interface{}
	Reply    chan string
	Function *luar.LuaObject
	Client   *Client
	Target   Targeter
	Line     *r1459.RawLine
}

// LoadScript finds and loads the appropriate script by a given short name (tetra/die).
func (tetra *Tetra) LoadScript(name string) (script *Script, err error) {
	kind := strings.Split(name, "/")[0]
	client, ok := tetra.Services[kind]
	if !ok {
		return nil, errors.New("Cannot find target service " + kind)
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
		Trigger:  make(chan []interface{}, 5),
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

	tetra.Etcd.CreateDir("/tetra/scripts/"+name, 0)

	go func() {
		for args := range script.Trigger {
			if len(args) == 2 {
				// Protocol hook
				debug("Protocol hook!")
				line, ok := args[1].(*r1459.RawLine)
				if !ok {
					debugf("Arg is %t, not *rfc1459.RawLine", args[1])
					return
				}
				debug(line.Raw)

				function, ok := args[0].(*luar.LuaObject)
				if !ok {
					debugf("Arg is %t, not *luar.LuaObject", args[0])
					return
				}
				debug(function.Type)

				res, err := function.Call(line)
				if err != nil {
					script.Log.Printf("Lua error %s: %s", script.Name, err.Error())
					script.Log.Printf("%#v", err)
					script.Client.ServicesLog(fmt.Sprintf("%s: %s", script.Name, err.Error()))
				}
				debug(res)

			} else {
				// Command
				debug("command")

				function, ok := args[0].(*luar.LuaObject)
				if !ok {
					debugf("Arg is %t, not *luar.LuaObject", args[0])
					return
				}

				client, ok := args[1].(*Client)
				if !ok {
					debugf("Arg is %t, not *Client", args[0])
					return
				}

				target, ok := args[2].(Targeter)
				if !ok {
					debugf("Arg is %t, not Targeter", args[0])
					return
				}

				cmdargs, ok := args[3].([]string)
				if !ok {
					debugf("Arg is %t, not []string", args[0])
					return
				}

				reschan, ok := args[4].(chan string)
				if !ok {
					debugf("Arg is %t, not chan string", args[0])
					return
				}

				reply, err := function.Call(client, target, cmdargs)
				if err != nil {
					script.Log.Printf("Lua error %s: %s", script.Name, err.Error())
					script.Log.Printf("%#v", err)
					script.Client.ServicesLog(fmt.Sprintf("%s: %s", script.Name, err.Error()))
					reschan <- ""
					return
				}

				reschan <- fmt.Sprintf("%s", reply)
			}
		}
	}()

	return
}

func (tetra *Tetra) loadLuaScript(script *Script) (*Script, error) {
	script.L.DoFile("modules/base.lua")

	err := script.L.DoFile("modules/" + script.Name + ".lua")

	if err != nil {
		return script, err
	}

	debugf("lua script %s loaded at %s", script.Name, script.Uuid)

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

	err := script.L.DoString(`moonscript = require "moonscript" xpcall = unsafe_xpcall pcall = unsafe_pcall local func, err = moonscript.loadstring(moonscript_code_from_file) if err ~= nil then tetra.log.Printf("Moonscript error, %#v", err) error(err) end func()`)
	if err != nil {
		script.Log.Print(err)
		return nil, err
	}

	debugf("moonscript script %s loaded at %s", script.Name, script.Uuid)

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
		"GC":        runtime.GC,
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

	luar.Register(script.L, "crypto", luar.Map{
		"hash": func(data string, salt string) string {
			output := md5.Sum([]byte(data + salt))
			return fmt.Sprintf("%x", output)
		},
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
		"format": func(format string, args ...interface{}) string {
			return fmt.Sprintf(format, args...)
		},
		"scan": fmt.Sscanf,
		"shuck": func(victim string) string {
			return victim[1 : len(victim)-1]
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

// AddLuaProtohook adds a lua function as a protocol hook.
func (script *Script) AddLuaProtohook(verb string, name string) error {
	function := luar.NewLuaObjectFromName(script.L, name)

	handler, err := script.Tetra.AddHandler(verb, func(line *r1459.RawLine) {
		debugf("sending %s", verb)
		script.Trigger <- []interface{}{function, line}
	})
	if err != nil {
		return err
	}

	handler.Script = script
	script.Handlers[handler.Uuid] = handler

	return nil
}

// AddLuaCommand adds a new command to a script from a lua context.
func (script *Script) AddLuaCommand(verb string, name string) error {
	function := luar.NewLuaObjectFromName(script.L, name)

	command, err := script.Client.NewCommand(verb, func(client *Client, target Targeter, args []string) string {
		reschan := make(chan string)
		defer close(reschan)

		script.Trigger <- []interface{}{
			function, client, target, args, reschan,
		}

		return <-reschan
	})

	if err != nil {
		return err
	}

	command.Script = script

	script.Commands[command.Uuid] = command

	return nil
}

// AddLuaHook adds a named hook from lua.
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
		delete(script.Client.Commands, command.Verb)
	}

	for _, hook := range script.Hooks {
		tetra.DelHook(hook)
	}

	script.L.Close()
	close(script.Trigger)

	delete(tetra.Scripts, name)

	return nil
}

func byteSliceToString(slice []byte) string {
	return string(slice)
}
