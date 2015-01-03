package tetra

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Xe/Tetra/1459"
	"github.com/Xe/Tetra/bot/script/charybdis"
	"github.com/Xe/Tetra/bot/script/crypto"
	tstrings "github.com/Xe/Tetra/bot/script/strings"
	lua "github.com/aarzilli/golua/lua"
	"github.com/sjkaliski/go-yo"
	"github.com/stevedonovan/luar"
)

// Struct Script implements a Lua scripting interface to Tetra.
type Script struct {
	Name     string
	L        *lua.State
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

var (
	Libraries map[string]luar.Map
)

func init() {
	Libraries = map[string]luar.Map{
		"crypto": luar.Map{
			"hash": crypto.Hash,
			"fnv":  crypto.Fnv,
		},

		"strings": luar.Map{
			"join":  strings.Join,
			"split": strings.Split,
			"first": tstrings.First,
			"rest":  tstrings.Rest,
			"format": func(format string, args ...interface{}) string {
				return fmt.Sprintf(format, args...)
			},
			"scan":  fmt.Sscanf,
			"shuck": tstrings.Shuck,
			"hassuffix": func(s, pattern string) bool {
				return strings.HasSuffix(s, pattern)
			},
		},

		"charybdis": luar.Map{
			"cloakhost": charybdis.CloakHost,
			"cloakip":   charybdis.CloakIP,
		},
	}
}

// LoadScript finds and loads the appropriate script by a given short name (tetra/die).
func LoadScript(name string) (script *Script, err error) {
	kind := strings.Split(name, "/")[0]
	client, ok := Services[kind]
	if !ok {
		return nil, errors.New("Cannot find target service " + kind)
	}

	if _, present := Scripts[name]; present {
		return nil, errors.New("Double script load!")
	}

	script = &Script{
		Name:     name,
		L:        luar.Init(),
		Log:      log.New(os.Stdout, name+" ", log.LstdFlags),
		Handlers: make(map[string]*Handler),
		Commands: make(map[string]*Command),
		Service:  kind,
		Client:   client,
		Uuid:     uuid.New(),
		Trigger:  make(chan []interface{}, 5),
	}

	script.seed()

	script, err = loadLuaScript(script)
	if err != nil {
		script, err = loadMoonScript(script)

		if err != nil {
			return nil, errors.New("No such script " + name)
		}
	}

	Scripts[name] = script

	Etcd.CreateDir("/tetra/scripts/"+name, 0)

	go func() {
		for args := range script.Trigger {
			switch args[0] {
			case INV_COMMAND:
				// Command
				debug("command")

				function, ok := args[1].(*luar.LuaObject)
				if !ok {
					debugf("Arg is %t, not *luar.LuaObject", args[0])
					return
				}

				client, ok := args[2].(*Client)
				if !ok {
					debugf("Arg is %t, not *Client", args[0])
					return
				}

				target, ok := args[3].(Targeter)
				if !ok {
					debugf("Arg is %t, not Targeter", args[0])
					return
				}

				cmdargs, ok := args[4].([]string)
				if !ok {
					debugf("Arg is %t, not []string", args[0])
					return
				}

				reschan, ok := args[5].(chan string)
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
			case INV_PROHOOK:
				// Protocol hook
				debug("Protocol hook!")
				line, ok := args[2].(*r1459.RawLine)
				if !ok {
					debugf("Arg is %t, not *rfc1459.RawLine", args[1])
					return
				}
				debug(line.Raw)

				function, ok := args[1].(*luar.LuaObject)
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
			case INV_NAMHOOK:
				debug("named hook")

				function, ok := args[1].(*luar.LuaObject)
				if !ok {
					debugf("Arg is %t, not *luar.LuaObject", args[1])
					return
				}
				debug(function.Type)

				funargs, ok := args[2].([]interface{})
				if !ok {
					debugf("Arg is %t, not []interface{}", args[2])
					return
				}

				_, err := function.Call(funargs...)
				if err != nil {
					script.Log.Printf("Lua error %s: %s", script.Name, err.Error())
					script.Log.Printf("%#v", err)
					script.Client.ServicesLog(fmt.Sprintf("%s: %s", script.Name, err.Error()))
				}
			}
		}
	}()

	return
}

func loadLuaScript(script *Script) (*Script, error) {
	script.L.DoFile("modules/base.lua")

	err := script.L.DoFile("modules/" + script.Name + ".lua")

	if err != nil {
		return script, err
	}

	debugf("lua script %s loaded at %s", script.Name, script.Uuid)

	script.Kind = "lua"

	return script, nil
}

func loadMoonScript(script *Script) (*Script, error) {
	contents, failed := ioutil.ReadFile("modules/" + script.Name + ".moon")

	if failed != nil {
		return script, errors.New("Could not read " + script.Name + ".moon")
	}

	luar.Register(script.L, "", luar.Map{
		"moonscript_code_from_file": string(contents),
	})

	/*
		moonscript = require "moonscript"
		xpcall = unsafe_xpcall
		pcall = unsafe_pcall
		local func, err = moonscript.loadstring(moonscript_code_from_file)
		if err ~= nil then
			tetra.log.Printf("Moonscript error, %#v", err)
			error(err)
		end
		func()
	*/
	err := script.L.DoString(`moonscript = require "moonscript" xpcall = unsafe_xpcall pcall = unsafe_pcall local func, err = moonscript.loadstring(moonscript_code_from_file) if err ~= nil then log.Printf("Moonscript error, %#v", err) error(err) end func()`)
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
		"use": func(library string) bool {
			if funcs, ok := Libraries[library]; ok {
				luar.Register(script.L, library, funcs)
				return true
			} else {
				return false
			}
		},
	})

	luar.Register(script.L, "tetra", luar.Map{
		"script":       script,
		"log":          script.Log,
		"Info":         Info,
		"Clients":      Clients,
		"Channels":     Channels,
		"Bursted":      Bursted,
		"Services":     Services,
		"Config":       ActiveConfig,
		"ActiveConfig": ActiveConfig,
		"Log":          script.Log,
		"Etcd":         Etcd,
		"Atheme":       Atheme,
		"RunHook":      RunHook,
		"LoadScript":   func(name string) (script *Script, err error) { return LoadScript(name) },
		"UnloadScript": func(name string) error { return UnloadScript(name) },
		"Scripts":      Scripts,
		"GetYo":        func(name string) (client *yo.Client, err error) { return GetYo(name) },
		"protohook":    script.AddLuaProtohook,
		"GC":           runtime.GC,
		"debug":        debug,
		"debugf":       debugf,
		"atheme":       Atheme,
		"Quit":         Quit,
	})

	luar.Register(script.L, "uuid", luar.Map{
		"new": uuid.New,
	})

	luar.Register(script.L, "web", luar.Map{
		"get":    http.Get,
		"post":   http.Post,
		"encode": url.QueryEscape,
		"decode": url.QueryUnescape,
	})

	luar.Register(script.L, "ioutil", luar.Map{
		"readall":     ioutil.ReadAll,
		"byte2string": byteSliceToString,
	})

	//luar.Register(script.L, "strings", Libraries["strings"])
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

	handler, err := AddHandler(verb, func(line *r1459.RawLine) {
		debugf("sending %s", verb)
		script.Trigger <- []interface{}{INV_PROHOOK, function, line}
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
			INV_COMMAND, function, client, target, args, reschan,
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

	hook := NewHook(verb, func(args ...interface{}) {
		script.Trigger <- []interface{}{
			INV_NAMHOOK, function, args,
		}
	})

	script.Hooks = append(script.Hooks, hook)

	return nil
}

// Unload a script and delete its commands and handlers
func UnloadScript(name string) error {
	if _, ok := Scripts[name]; !ok {
		panic("No such script " + name)
	}

	script := Scripts[name]

	for _, handler := range script.Handlers {
		DelHandler(handler.Verb, handler.Uuid)
		delete(script.Handlers, handler.Uuid)
	}

	for _, command := range script.Commands {
		delete(script.Commands, command.Uuid)
		delete(script.Client.Commands, command.Verb)
	}

	for _, hook := range script.Hooks {
		if hook.Verb == "SHUTDOWN" {
			hook.impl()
		}

		DelHook(hook)
	}

	script.L.Close()
	close(script.Trigger)

	delete(Scripts, name)

	return nil
}

func byteSliceToString(slice []byte) string {
	return string(slice)
}
