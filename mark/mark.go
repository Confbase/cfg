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
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if cfg.UnMark {
		unmark.Unmark(cfg.Targets)
		os.Exit(0)
	}
	if cfg.Singleton {
		for _, target := range cfg.Targets {
			track.Track(baseDir, target)
		}
		os.Exit(0)
	}
	if cfg.InstanceOf != "" {
		for _, target := range cfg.Targets {
			tag.MustTag(target, cfg.InstanceOf)
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
	absPath, relPath, err := dotcfg.GetAbsAndRelPaths(baseDir, cfg.Targets[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	cfgFile := dotcfg.MustLoadCfg(baseDir)

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
		FilePath: relPath,
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
		cfgFile.Templates[templIndex] = templObj
	}
	if err := cfgFile.Infer(baseDir, cfg.Targets[0]); err == nil {
		// if infer was successful
		for _, inst := range cfgFile.Instances {
			for _, templName := range inst.TemplNames {
				if templName == templObj.Name {
					cfgFile.MustWarnDiffs(baseDir, templObj.Name, inst.FilePath, os.Stderr)
					break
				}
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "warning: failed to infer schema of %v\n", absPath)
	}
	cfgFile.MustSerialize(baseDir, nil)
	if !cfgFile.NoGit {
		cfgFile.MustStage(baseDir)
		cfgFile.MustCommit(baseDir)
	}
}
