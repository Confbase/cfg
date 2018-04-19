package track

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/lib/dotcfg"
)

func Track(filePath string) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: the file '%v' does not exist\n", filePath)
		os.Exit(1)
	}

	cfg := dotcfg.MustLoadCfg()

	containsSingleton := false
	for _, singleton := range cfg.Singletons {
		if singleton == filePath {
			containsSingleton = true
			break
		}
	}

	if containsSingleton {
		fmt.Fprintf(os.Stderr, "error: '%v' is already tracked as an singleton\n", filePath)
		os.Exit(1)
	}

	cfg.Singletons = append(cfg.Singletons, filePath)
	cfg.MustSerialize(nil)
	cfg.MustStage()
	cfg.MustCommit()
}
