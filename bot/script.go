package tetra

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/Xe/Tetra/1459"
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"errors"
	"net/http"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Script struct {
	Name     string
	L        *lua.State
	Tetra    *Tetra
	Log      *log.Logger
	Handlers map[string]*Handler
	Service  string
	Client   *Client
	Uuid     string
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
		Service:  kind,
		Client:   client,
		Uuid:     uuid.New(),
	}

	luar.Register(script.L, "", luar.Map{
		"client": script.Client,
	})

	luar.Register(script.L, "tetra", luar.Map{
		"script": script,
		"log":    script.Log,
		"bot":    tetra,
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

	script.L.DoFile("modules/base.lua")

	tetra.Scripts[name] = script

	err = script.L.DoFile("modules/" + name + ".lua")
	if err != nil {
		return nil, err
	}

	return
}

// Add a lua function as a protocol hook
func (script *Script) AddLuaProtohook(verb string, name string) {
	function := luar.NewLuaObjectFromName(script.L, name)

	handler, err := script.Tetra.AddHandler(verb, func(line *r1459.RawLine) {
		_, err := function.Call(line)
		if err != nil {
			script.Log.Printf("Lua error %s: %#v", script.Name, err)
		}
	})
	if err != nil {
		panic(err)
	}

	handler.Script = script
	script.Handlers[handler.Uuid] = handler
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

	script.L.Close()

	delete(tetra.Scripts, name)

	return nil
}

func byteSliceToString(slice []byte) string {
	return string(slice)
}
