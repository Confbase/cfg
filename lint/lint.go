package lint

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/Confbase/cfg/dotcfg"
)

func MustLint() {
	if err := Lint(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Lint() error {
	cfgFile, err := dotcfg.LoadCfg("")
	if err != nil {
		return err
	}
	for _, templ := range cfgFile.Templates {
		// TODO: add config variable + data structure to only do this to certain files
		if templ.Schema.FilePath == "" {
			if err := cfgFile.Infer(templ.FilePath); err != nil {
				exiterr, ok := err.(*exec.ExitError)
				if !ok {
					return err
				}
				// 'schema diff' exited with an exit code != 0
				status, ok := exiterr.Sys().(syscall.WaitStatus)
				if !ok {
					return err
				}
				if status.ExitStatus() == 1 || status.ExitStatus() == 127 {
					// TODO: add config variable to suppress these warnings
					fmt.Printf("warning: cannot lint '%v' template\n", templ.Name)
					continue
				}
				return err
			}
		}
		for _, inst := range cfgFile.Instances {
			isInstOfTempl := false
			for _, templName := range inst.TemplNames {
				if templName == templ.Name {
					isInstOfTempl = true
					break
				}
			}
			if isInstOfTempl {
				// TODO: add config variable + data structure to only do this to certain files
				if inst.FilePath == "" {
					if err := cfgFile.Infer(inst.FilePath); err != nil {
						return err
					}
				}
				if err := cfgFile.WarnDiffs(templ.Name, inst.FilePath, os.Stderr); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
