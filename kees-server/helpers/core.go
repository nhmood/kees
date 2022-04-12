package helpers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"kees-server/config"
	"runtime"
)

var debug bool

func Configure(config *config.Config) {
	debug = config.Logging.Debug
	spew.Config.MaxDepth = config.Logging.SpewDepth
}

func Dump(data interface{}) {
	spew.Dump(data)
}

func Debug(data interface{}) {
	if debug {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("called from %s#%d\n", file, no)
		}
		Dump(data)
		fmt.Printf("===================================\n\n\n")
	}
}
