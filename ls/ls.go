package ls

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Confbase/cfg/decorate"
	"github.com/Confbase/cfg/dotcfg"
)

func isStdoutTty() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func Ls(lsCfg *Config) {
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cfg := dotcfg.MustLoadCfg(baseDir)

	if isStdoutTty() && !lsCfg.NoTty {
		d := decorate.New()
		d.Enabled = !lsCfg.NoColors

		snaps := dotcfg.MustLoadSnaps(baseDir)
		fmt.Printf("## %v\n", snaps.Current.Name)

		if !(lsCfg.DoLsTempls || lsCfg.DoLsInsts || lsCfg.DoLsSingles || lsCfg.DoLsUntracked) {
			LsTemplsHuman(cfg, d)
			fmt.Println()
			LsInstsHuman(cfg, d)
			fmt.Println()
			LsSinglesHuman(cfg, d)
			return
		}
		if lsCfg.DoLsTempls {
			LsTemplsHuman(cfg, d)
			if lsCfg.DoLsInsts || lsCfg.DoLsSingles || lsCfg.DoLsUntracked {
				fmt.Println()
			}
		}
		if lsCfg.DoLsInsts {
			LsInstsHuman(cfg, d)
			if lsCfg.DoLsSingles || lsCfg.DoLsUntracked {
				fmt.Println()
			}
		}
		if lsCfg.DoLsSingles {
			LsSinglesHuman(cfg, d)
			if lsCfg.DoLsUntracked {
				fmt.Println()
			}
		}
		if lsCfg.DoLsUntracked {
			LsUntrackedHuman(cfg, d)
		}
	} else {
		if !(lsCfg.DoLsTempls || lsCfg.DoLsInsts || lsCfg.DoLsSingles || lsCfg.DoLsUntracked) {
			LsTemplsTty(cfg)
			LsInstsTty(cfg)
			LsSinglesTty(cfg)
			return
		}
		if lsCfg.DoLsTempls {
			LsTemplsTty(cfg)
		}
		if lsCfg.DoLsInsts {
			LsInstsTty(cfg)
		}
		if lsCfg.DoLsSingles {
			LsSinglesTty(cfg)
		}
		if lsCfg.DoLsUntracked {
			LsUntrackedTty(cfg)
		}
	}
}

func LsTemplsHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("templates")))
	if len(cfg.Templates) > 0 {
		for _, t := range cfg.Templates {
			relativePath, err := dotcfg.ConvertPathToRelative(filepath.Join(cfg.BaseDir, t.FilePath))
			if err != nil {
				// This seems somewhat hacky... but I'm not sure how best to handle this
				relativePath = t.FilePath
			}

			fmt.Printf("%v: %v\n", d.Green(t.Name), relativePath)
		}
	}
}

func LsTemplsTty(cfg *dotcfg.File) {
	fmt.Println("templates")
	for _, t := range cfg.Templates {
		relativePath, err := dotcfg.ConvertPathToRelative(filepath.Join(cfg.BaseDir, t.FilePath))
		if err != nil {
			// This seems somewhat hacky... but I'm not sure how best to handle this
			relativePath = t.FilePath
		}

		fmt.Printf("%v\t%v\n", t.Name, relativePath)
	}
}

func LsInstsHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("instances")))
	if len(cfg.Templates) > 0 {
		for _, inst := range cfg.Instances {
			templsStr := strings.Join(inst.TemplNames, ", ")
			relativePath, err := dotcfg.ConvertPathToRelative(filepath.Join(cfg.BaseDir, inst.FilePath))
			if err != nil {
				// This seems somewhat hacky... but I'm not sure how best to handle this
				relativePath = inst.FilePath
			}

			fmt.Printf("%v: %v\n", d.Green(relativePath), templsStr)
		}
	}
}

func LsInstsTty(cfg *dotcfg.File) {
	fmt.Println("instances")
	for _, inst := range cfg.Instances {
		templsStr := strings.Join(inst.TemplNames, ",")
		relativePath, err := dotcfg.ConvertPathToRelative(filepath.Join(cfg.BaseDir, inst.FilePath))
		if err != nil {
			// This seems somewhat hacky... but I'm not sure how best to handle this
			relativePath = inst.FilePath
		}

		fmt.Printf("%v\t%v\n", relativePath, templsStr)
	}
}

func LsSinglesHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("singletons")))
	for _, s := range cfg.Singletons {
		relativePath, err := dotcfg.ConvertPathToRelative(filepath.Join(cfg.BaseDir, s.FilePath))
		if err != nil {
			// This seems somewhat hacky... but I'm not sure how best to handle this
			relativePath = s.FilePath
		}

		fmt.Printf("%v\n", relativePath)
	}
}

func LsSinglesTty(cfg *dotcfg.File) {
	fmt.Println("singletons")
	for _, s := range cfg.Singletons {
		relativePath, err := dotcfg.ConvertPathToRelative(filepath.Join(cfg.BaseDir, s.FilePath))
		if err != nil {
			// This seems somewhat hacky... but I'm not sure how best to handle this
			relativePath = s.FilePath
		}

		fmt.Println(relativePath)
	}
}

func LsUntrackedHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	untrackedFiles, err := cfg.GetUntrackedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(d.LightBlue(d.Title("untracked files")))
	for _, uf := range untrackedFiles {
		fmt.Println(uf)
	}
}

func LsUntrackedTty(cfg *dotcfg.File) {
	untrackedFiles, err := cfg.GetUntrackedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("untracked files")
	for _, uf := range untrackedFiles {
		fmt.Println(uf)
	}
}
