package track

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
)

// TODO: MustTrack

// Track tracks the file located at the absolute or relative path `filePath`.
// `baseDir` specifies the base directory of this base. It has no relation to
// and does not modify `filePath`.
func Track(baseDir, filePath string) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: the file '%v' does not exist\n", filePath)
		os.Exit(1)
	}

	cfg := dotcfg.MustLoadCfg(baseDir)

	containsSingleton := false
	for _, singleton := range cfg.Singletons {
		if singleton.FilePath == filePath {
			containsSingleton = true
			break
		}
	}

	if containsSingleton {
		fmt.Fprintf(os.Stderr, "error: '%v' is already tracked as an singleton\n", filePath)
		os.Exit(1)
	}

	cfg.Singletons = append(cfg.Singletons, dotcfg.Singleton{FilePath: filePath})
	cfg.Infer(filePath)
	cfg.MustSerialize(baseDir, nil)
	if !cfg.NoGit {
		cfg.MustStage(baseDir)
		cfg.MustCommit(baseDir)
	}
}
