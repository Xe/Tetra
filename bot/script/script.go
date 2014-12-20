package script

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
	"github.com/Xe/Tetra/bot/script/crypto"
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

// Struct Script implements a Lua scripting interface to Tetra.
type Script struct {
	Name    string
	L       *lua.State
	Log     *log.Logger
	Service string
	Uuid    string
	Kind    string
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
	Line     *r1459.RawLine
}

// LoadScript finds and loads the appropriate script by a given short name (tetra/die).
func NewScript(name string) (script *Script, err error) {
	kind := strings.Split(name, "/")[0]

	script = &Script{
		Name:    name,
		L:       luar.Init(),
		Log:     log.New(os.Stdout, name+" ", log.LstdFlags),
		Service: kind,
		Uuid:    uuid.New(),
	}

	script.seed()

	script, err = loadLuaScript(script)
	if err != nil {
		script, err = loadMoonScript(script)

		if err != nil {
			return nil, errors.New("No such script " + name)
		}
	}

	return
}

func loadLuaScript(script *Script) (*Script, error) {
	script.L.DoFile("modules/base.lua")

	err := script.L.DoFile("modules/" + script.Name + ".lua")

	if err != nil {
		return script, err
	}

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

	err := script.L.DoString(`moonscript = require "moonscript" xpcall = unsafe_xpcall pcall = unsafe_pcall local func, err = moonscript.loadstring(moonscript_code_from_file) if err ~= nil then tetra.log.Printf("Moonscript error, %#v", err) error(err) end func()`)
	if err != nil {
		script.Log.Print(err)
		return nil, err
	}

	script.Kind = "moonscript"

	return script, nil
}

func (script *Script) seed() {
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
		"hash": crypto.Hash,
		"fnv":  crypto.Fnv,
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
		"hassuffix": func(s, pattern string) bool {
			return strings.HasSuffix(s, pattern)
		},
	})
}

// Unload a script and delete its commands and handlers
func (s *Script) Unload(name string) error {
	s.L.Close()

	return nil
}

func byteSliceToString(slice []byte) string {
	return string(slice)
}
