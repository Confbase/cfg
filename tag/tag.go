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

	cfgFile := dotcfg.MustLoadCfg()

	containsTempl := false
	for _, t := range cfgFile.Templates {
		if t.Name == templName {
			containsTempl = true
			break
		}
	}
	if !containsTempl {
		fmt.Fprintf(os.Stderr, "error: template '%v' does not exist\n", templName)

		templsContainPath := false
		guessTemplName := ""
		for _, t := range cfgFile.Templates {
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

	isNewInst := true
	for i, inst := range cfgFile.Instances {
		if inst.FilePath == filePath {
			for _, t := range cfgFile.Instances[i].TemplNames {
				if t == templName {
					// if already tagged as this templ,
					// do nothing
					return
				}
			}
			cfgFile.Instances[i].TemplNames = append(inst.TemplNames, templName)
			isNewInst = false
			break
		}
	}
	if isNewInst {
		inst := dotcfg.NewInstance(filePath)
		inst.TemplNames = append(inst.TemplNames, templName)
		cfgFile.Instances = append(cfgFile.Instances, *inst)
		if err := cfgFile.Infer(filePath); err == nil {
			// if infer was successful
			cfgFile.MustWarnDiffs(templName, filePath)
		}
	}

	// if target is already a singleton, remove it from the singletons list
	for i, s := range cfgFile.Singletons {
		if s.FilePath == filePath {
			cfgFile.Singletons = append(cfgFile.Singletons[:i], cfgFile.Singletons[i+1:]...)
			break
		}
	}

	cfgFile.MustSerialize(nil)
	if !cfgFile.NoGit {
		cfgFile.MustStage()
		cfgFile.MustCommit()
	}
}
