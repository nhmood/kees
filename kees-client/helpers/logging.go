package helpers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"kees/media-controller/config"
	"runtime"
)

var debug bool

func ConfigureLogging(config config.LoggingConfig) {
	debug = config.Debug
	spew.Config.MaxDepth = config.SpewDepth
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
