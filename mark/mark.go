package mark

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
	"github.com/Confbase/cfg/tag"
	"github.com/Confbase/cfg/track"
	"github.com/Confbase/cfg/unmark"
)

func Mark(cfg *Config) {
	if cfg.UnMark {
		unmark.Unmark(cfg.Targets)
		os.Exit(0)
	}
	if cfg.Singleton {
		for _, target := range cfg.Targets {
			track.Track(target)
		}
		os.Exit(0)
	}
	if cfg.InstanceOf != "" {
		for _, target := range cfg.Targets {
			tag.Tag(target, cfg.InstanceOf)
		}
		os.Exit(0)
	}
	if cfg.Template == "" {
		fmt.Fprintf(os.Stderr, "error: one of the flags (-u|-i|-t) is required; see 'cfg mark -h' for help\n")
		os.Exit(1)
	}

	if len(cfg.Targets) > 1 {
		fmt.Fprintf(os.Stderr, "any given template can only be associated to one file\n")
		os.Exit(1)
	}
	target := cfg.Targets[0]
	_, err := os.Stat(target)
	if err != nil && os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: the file '%v' does not exist\n", target)
		os.Exit(1)
	}

	cfgFile := dotcfg.MustLoadCfg()

	containsTempl := false
	templIndex := -1
	for i, t := range cfgFile.Templates {
		if t.Name == cfg.Template {
			containsTempl = true
			templIndex = i
			break
		}
	}

	templObj := dotcfg.Template{
		Name:     cfg.Template,
		FilePath: target,
	}

	if !containsTempl {
		cfgFile.Templates = append(
			cfgFile.Templates,
			templObj,
		)
	} else {
		if !cfg.Force {
			fmt.Fprintf(os.Stderr, "template '%v' already exists; ", cfg.Template)
			fmt.Fprintf(os.Stderr, "use --force to overwrite it\n")
			os.Exit(1)
		}

		oldTemplName := cfgFile.Templates[templIndex].Name
		cfgFile.Templates[templIndex] = templObj

		for i, _ := range cfgFile.Instances {
			tns := cfgFile.Instances[i].TemplNames
			for j, t := range tns {
				if t == oldTemplName {
					tns = append(tns[:j], tns[j+1:]...)
					break
				}
			}
			cfgFile.Instances[i].TemplNames = tns
		}
	}

	cfgFile.MustSerialize(nil)
	if !cfgFile.NoGit {
		cfgFile.MustStage()
		cfgFile.MustCommit()
	}
}
