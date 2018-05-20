package mark

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
)

func Mark(filePath, templ string, force bool) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: the file '%v' does not exist\n", filePath)
		os.Exit(1)
	}

	cfg := dotcfg.MustLoadCfg()

	templObj := dotcfg.Template{
		Name:     templ,
		FilePath: filePath,
	}

	containsTempl := false
	templIndex := -1
	for i, t := range cfg.Templates {
		if t.Name == templ {
			containsTempl = true
			templIndex = i
			break
		}
	}

	if !containsTempl {
		cfg.Templates = append(
			cfg.Templates,
			templObj,
		)
		fmt.Printf("created new template '%v'\n", templ)
	} else {
		if !force {
			fmt.Fprintf(os.Stderr, "template '%v' already exists; use --force to overwrite it\n", templ)
			os.Exit(1)
		}

		oldTemplName := cfg.Templates[templIndex].Name
		cfg.Templates[templIndex] = templObj

		if _, ok := cfg.Instances[oldTemplName]; ok {
			delete(cfg.Instances, oldTemplName)
		}
		fmt.Printf("overwrote template '%v'\n", templ)
	}

	cfg.MustSerialize(nil)
	if !cfg.NoGit {
		cfg.MustStage()
		cfg.MustCommit()
	}
}
