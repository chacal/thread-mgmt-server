package server

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

func ParseOptions(opts interface{}) {
	var parser = flags.NewParser(opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type != flags.ErrHelp {
			fmt.Println()
			parser.WriteHelp(os.Stdout)
		}
		os.Exit(1)
	}
}
