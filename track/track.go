package track

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
)

// TODO: MustTrack

// Track tracks the file located at the absolute or relative path `filePath`.
func Track(baseDir, filePath string) {
	absPath, relPath, err := dotcfg.GetAbsAndRelPaths(baseDir, filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	_, err = os.Stat(absPath)
	if err != nil && os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: the file '%v' does not exist\n", absPath)
		os.Exit(1)
	}

	cfg := dotcfg.MustLoadCfg(baseDir)

	containsSingleton := false
	for _, singleton := range cfg.Singletons {
		if singleton.FilePath == relPath {
			containsSingleton = true
			break
		}
	}

	if containsSingleton {
		fmt.Fprintf(os.Stderr, "error: '%v' is already tracked as an singleton\n", absPath)
		os.Exit(1)
	}

	cfg.Singletons = append(cfg.Singletons, dotcfg.Singleton{FilePath: relPath})
	cfg.Infer(baseDir, filePath)
	cfg.MustSerialize(baseDir, nil)
	if !cfg.NoGit {
		cfg.MustStage(baseDir)
		cfg.MustCommit(baseDir)
	}
}
