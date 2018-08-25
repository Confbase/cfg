package config

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
)

func Config(args []string) {
	keyFile := dotcfg.MustLoadKey()
	// len(args) guaranteed to be even
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]
		switch key {
		case "email":
			keyFile.Email = value
		case "base":
			keyFile.BaseName = value
		default:
			fmt.Fprintf(os.Stderr, "error: invalid key '%v'\n", key)
			os.Exit(1)
		}
	}
	keyFile.MustSerialize("", nil)
}
