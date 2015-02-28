package tetra

import (
	"runtime"
)

func AddHacks() {
	NewHook("SCRIPTLOAD", func(i ...interface{}) {
		runtime.GOMAXPROCS(len(Scripts))
	})

	NewHook("SCRIPTUNLOAD", func(i ...interface{}) {
		runtime.GOMAXPROCS(len(Scripts) - 1)
	})
}
