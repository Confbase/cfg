package unmark

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
)

func MustUnmark(targets []string) {
	if err := Unmark(targets); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Unmark(targets []string) error {
	cfgFile, err := dotcfg.LoadCfg("")
	if err != nil {
		return err
	}
	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			return fmt.Errorf("failed to stat '%v'\n%v", target, err)
		}

		if !cfgFile.NoGit {
			cmd := exec.Command("git", "ls-files", "--error-unmatch", target)
			if err := cmdrunner.PipeTo(cmd, nil, nil); err != nil {
				return fmt.Errorf("'%v' is not tracked by git\n%v", target, err)
			}
		}

		var rmdTempl dotcfg.Template
		isTemplate := false

		for i, t := range cfgFile.Templates {
			if t.FilePath == target {
				cfgFile.Templates = append(cfgFile.Templates[:i], cfgFile.Templates[i+1:]...)
				rmdTempl = t
				isTemplate = true
				if err := cfgFile.RmSchema(target, rmdTempl.Schema.FilePath != ""); err != nil {
					return err
				}
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
				} else {
					if err := cfgFile.RmSchema(inst.FilePath, inst.Schema.FilePath != ""); err != nil {
						return err
					}
				}
			}
			cfgFile.Instances = insts
		}

		insts := cfgFile.Instances
		for i, inst := range insts {
			if inst.FilePath == target {
				insts = append(insts[:i], insts[i+1:]...)
				if err := cfgFile.RmSchema(target, inst.Schema.FilePath != ""); err != nil {
					return err
				}
				break
			}
		}
		cfgFile.Instances = insts

		ss := cfgFile.Singletons
		for i, s := range ss {
			if s.FilePath == target {
				ss = append(ss[:i], ss[i+1:]...)
				if err := cfgFile.RmSchema(target, s.Schema.FilePath != ""); err != nil {
					return err
				}
				break
			}
		}
		cfgFile.Singletons = ss

		if !cfgFile.NoGit {
			out, err := exec.Command("git", "rm", "--cached", target).CombinedOutput()
			if err != nil {
				return fmt.Errorf("'git rm --cached %v' failed\nerror: %v\noutput: %v", target, err, string(out))
			}
		}
	}

	if err := cfgFile.Serialize("", nil); err != nil {
		return err
	}
	if !cfgFile.NoGit {
		if err := cfgFile.Stage(""); err != nil {
			return err
		}
		if err := cfgFile.Commit(""); err != nil {
			return err
		}
	}
	return nil
}
