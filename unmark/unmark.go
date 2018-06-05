package unmark

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
)

func Unmark(targets []string) {
	cfgFile := dotcfg.MustLoadCfg()
	for _, target := range targets {

		if _, err := os.Stat(target); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to stat '%v'\n", target)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		if !cfgFile.NoGit {
			cmd := exec.Command("git", "ls-files", "--error-unmatch", target)
			if err := cmdrunner.PipeFrom(cmd, nil, nil); err != nil {
				fmt.Fprintf(os.Stderr, "error: '%v' is not tracked by git\n", target)
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		}

		var rmdTempl dotcfg.Template
		isTemplate := false

		for i, t := range cfgFile.Templates {
			if t.FilePath == target {
				cfgFile.Templates = append(cfgFile.Templates[:i], cfgFile.Templates[i+1:]...)
				rmdTempl = t
				isTemplate = true
				break
			}
		}
		if isTemplate {
			insts := make([]dotcfg.Instance, 0)
			for _, inst := range cfgFile.Instances {
				for i, t := range inst.TemplNames {
					if t == rmdTempl.Name {
						inst.TemplNames = append(inst.TemplNames[:i], inst.TemplNames[i+1:]...)
						break
					}
				}
				if len(inst.TemplNames) != 0 {
					insts = append(insts, inst)
				}
			}
			cfgFile.Instances = insts
		}

		insts := cfgFile.Instances
		for i, inst := range insts {
			if inst.FilePath == target {
				insts = append(insts[:i], insts[i+1:]...)
				break
			}
		}
		cfgFile.Instances = insts

		ss := cfgFile.Singletons
		for i, s := range ss {
			if s.FilePath == target {
				ss = append(ss[:i], ss[i+1:]...)
				break
			}
		}
		cfgFile.Singletons = ss

		if !cfgFile.NoGit {
			out, err := exec.Command("git", "rm", "--cached", target).CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "'git rm --cached %v' failed\n", target)
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				fmt.Fprintf(os.Stderr, "output: %v\n", string(out))
				os.Exit(1)
			}
		}
	}

	cfgFile.MustSerialize(nil)
	if !cfgFile.NoGit {
		cfgFile.MustStage()
		cfgFile.MustCommit()
	}
}
