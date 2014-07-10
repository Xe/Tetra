package tetra

import (
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"log"
	"os"
)

type Script struct {
	Name     string
	L        *lua.State
	Tetra    *Tetra
	Log      *log.Logger
	Commands []*Command
	Handlers []*Handler
}

func (tetra *Tetra) LoadScript(name string) (script *Script) {
	script = &Script{
		Name:     name,
		L:        luar.Init(),
		Tetra:    tetra,
		Log:      log.New(os.Stdout, name+" ", log.LstdFlags),
		Commands: nil,
		Handlers: nil,
	}

	luar.Register(script.L, "tetra", luar.Map{
		"script": script,
		"log":    script.Log,
		"bot":    tetra,
	})

	tetra.Scripts[name] = script

	script.L.DoFile(name)

	return
}
