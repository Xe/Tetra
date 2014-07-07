package cod

import (
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"log"
)

type Script struct {
	L        *lua.State
	Cod      *Cod
	Log      *log.Logger
	Commands []*Command
	Handlers []*Handler
}

func (s *Script) Register() {
	luar.Register(s.L, "cod", luar.Map{
		"script": s,
		"log": s.Log,
	})
}

