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
			fmt.Fprintf(os.Stderr, "however, the file '%v' is ", templName)
			fmt.Fprintf(os.Stderr, "associated with the name '%v'\n", guessTemplName)
			fmt.Fprintf(os.Stderr, "did you mean to run ")
			fmt.Fprintf(os.Stderr, "'cfg mark -i %v %v'?\n", guessTemplName, filePath)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "use 'cfg mark -t' to mark a file as a template ")
		fmt.Fprintf(os.Stderr, "before marking an instance of it\n")
		os.Exit(1)
	}

	if _, ok := cfg.Instances[templName]; !ok {
		cfg.Instances[templName] = make([]dotcfg.Instance, 0)
	}

	containsInstance := false
	for _, instance := range cfg.Instances[templName] {
		if instance.FilePath == filePath {
			containsInstance = true
			break
		}
	}

	if containsInstance {
		fmt.Fprintf(os.Stderr, "error: '%v' is already tagged as an instance of '%v'\n", filePath, templName)
		os.Exit(1)
	}

	cfg.Instances[templName] = append(cfg.Instances[templName], dotcfg.Instance{FilePath: filePath})

	cfg.MustSerialize(nil)
	if !cfg.NoGit {
		cfg.MustStage()
		cfg.MustCommit()
	}
}
