package tag

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
)

func MustTag(filePath, templName string) {
	if err := Tag(filePath, templName); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Tag(filePath, templName string) error {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	cfgFile, err := dotcfg.LoadCfg("")
	if err != nil {
		return err
	}

	containsTempl := false
	for _, t := range cfgFile.Templates {
		if t.Name == templName {
			containsTempl = true
			break
		}
	}
	if !containsTempl {
		err := fmt.Errorf("template '%v' does not exist\n", templName)

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
			err = fmt.Errorf("%vhowever, the file '%v' is ", err, templName)
			err = fmt.Errorf("%vassociated with the name '%v'\n", err, guessTemplName)
			err = fmt.Errorf("%vdid you mean to run ", err)
			return fmt.Errorf("%v'cfg mark -i %v %v'?", err, guessTemplName, filePath)
		}

		err = fmt.Errorf("%vuse 'cfg mark -t' to mark a file as a template ", err)
		return fmt.Errorf("%vbefore marking an instance of it", err)
	}

	isNewInst := true
	for i, inst := range cfgFile.Instances {
		if inst.FilePath == filePath {
			for _, t := range cfgFile.Instances[i].TemplNames {
				if t == templName {
					// if already tagged as this templ,
					// do nothing
					return nil
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
			if err := cfgFile.WarnDiffs(templName, filePath, os.Stderr); err != nil {
				return err
			}
		}
	}

	// if target is already a singleton, remove it from the singletons list
	for i, s := range cfgFile.Singletons {
		if s.FilePath == filePath {
			cfgFile.Singletons = append(cfgFile.Singletons[:i], cfgFile.Singletons[i+1:]...)
			break
		}
	}

	if err := cfgFile.Serialize("", nil); err != nil {
		return err
	}
	if !cfgFile.NoGit {
		if err := cfgFile.Stage(""); err != nil {
			return err
		}
		return cfgFile.Commit("")
	}
	return nil
}
