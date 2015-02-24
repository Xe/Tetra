package tetra

import (
	"runtime"
)

func init() {
	NewHook("SCRIPTLOAD", func(i ...interface{}) {
		runtime.GOMAXPROCS(len(Scripts))
	})

	NewHook("SCRIPTUNLOAD", func(i ...interface{}) {
		runtime.GOMAXPROCS(len(Scripts) - 1)
	})
}
