package tag

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
)

func Tag(filePath, templName string) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: the file '%v' does not exist\n", filePath)
		os.Exit(1)
	}

	cfg := dotcfg.MustLoadCfg()

	containsTempl := false
	for _, t := range cfg.Templates {
		if t.Name == templName {
			containsTempl = true
			break
		}
	}
	if !containsTempl {
		fmt.Fprintf(os.Stderr, "error: template '%v' does not exist\n", templName)

		templsContainPath := false
		guessTemplName := ""
		for _, t := range cfg.Templates {
			if t.FilePath == templName {
				templsContainPath = true
				guessTemplName = t.Name
				break
			}
		}

		if templsContainPath {
			fmt.Fprintf(
				os.Stderr,
				"however, the file '%v' is associated with the name '%v'\ndid you mean to run 'cfg mark %v %v' instead?\n",
				templName,
				guessTemplName,
				filePath,
				guessTemplName,
			)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "use 'cfg mark' to specify which file is the template before tagging an instance of it\n")
		os.Exit(1)
	}

	if _, ok := cfg.Instances[templName]; !ok {
		cfg.Instances[templName] = make([]string, 0)
	}

	containsInstance := false
	for _, instance := range cfg.Instances[templName] {
		if instance == filePath {
			containsInstance = true
			break
		}
	}

	if containsInstance {
		fmt.Fprintf(os.Stderr, "error: '%v' is already tagged as an instance of '%v'\n", filePath, templName)
		os.Exit(1)
	}

	cfg.Instances[templName] = append(cfg.Instances[templName], filePath)
	fmt.Printf("tagged '%v' as an instance of '%v'\n", filePath, templName)

	cfg.MustSerialize(nil)
	if !cfg.NoGit {
		cfg.MustStage()
		cfg.MustCommit()
	}
}
