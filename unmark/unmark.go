package unmark

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Confbase/cfg/dotcfg"
)

func Unmark(targets []string) {
	cfgFile := dotcfg.MustLoadCfg()
	for _, target := range targets {

		var rmdTempl dotcfg.Template
		isTemplate := false

		ts := cfgFile.Templates
		for i, t := range ts {
			if t.FilePath == target {
				if i == len(ts)-1 {
					ts = ts[:len(ts)-1]
				} else {
					ts = append(ts[:i], ts[i+1:]...)
				}
				rmdTempl = t
				isTemplate = true
				break
			}
		}
		cfgFile.Templates = ts
		if isTemplate {
			for templName, insts := range cfgFile.Instances {
				if len(insts) == 0 {
					continue
				}
				if templName == rmdTempl.Name {
					fmt.Fprintf(os.Stderr, "error: cannot unmark ")
					fmt.Fprintf(os.Stderr, "'%v'\nit is the ", target)
					fmt.Fprintf(os.Stderr, "template '%v' ", templName)
					fmt.Fprintf(os.Stderr, "and there are instances of it")
					os.Exit(1)
				}
			}
		}

		for i, templ := range cfgFile.Instances {
			for i, inst := range templ {
				if inst.FilePath == target {
					if i == len(templ)-1 {
						templ = templ[:len(templ)-1]
					} else {
						templ = append(templ[:i], templ[i+1:]...)
					}
					break
				}
			}
			cfgFile.Instances[i] = templ
		}

		ss := cfgFile.Singletons
		for i, s := range ss {
			if s.FilePath == target {
				if i == len(ss)-1 {
					ss = ss[:len(ss)-1]
				} else {
					ss = append(ss[:i], ss[i+1:]...)
				}
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
