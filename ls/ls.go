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
		baseDir, err := dotcfg.GetBaseDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		for _, t := range cfg.Templates {
			absPath := filepath.Join(baseDir, t.FilePath)
			relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%v: %v\n", d.Green(t.Name), relPath)
		}
	}
}

func LsTemplsTty(cfg *dotcfg.File) {
	fmt.Println("templates")
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, t := range cfg.Templates {
		absPath := filepath.Join(baseDir, t.FilePath)
		relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%v\t%v\n", t.Name, relPath)
	}
}

func LsInstsHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("instances")))
	if len(cfg.Templates) > 0 {
		baseDir, err := dotcfg.GetBaseDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		for _, inst := range cfg.Instances {
			absPath := filepath.Join(baseDir, inst.FilePath)
			relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			templsStr := strings.Join(inst.TemplNames, ", ")
			fmt.Printf("%v: %v\n", d.Green(relPath), templsStr)
		}
	}
}

func LsInstsTty(cfg *dotcfg.File) {
	fmt.Println("instances")
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, inst := range cfg.Instances {
		absPath := filepath.Join(baseDir, inst.FilePath)
		relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		templsStr := strings.Join(inst.TemplNames, ",")
		fmt.Printf("%v\t%v\n", relPath, templsStr)
	}
}

func LsSinglesHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("singletons")))
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, s := range cfg.Singletons {
		absPath := filepath.Join(baseDir, s.FilePath)
		relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%v\n", relPath)
	}
}

func LsSinglesTty(cfg *dotcfg.File) {
	fmt.Println("singletons")
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, s := range cfg.Singletons {
		absPath := filepath.Join(baseDir, s.FilePath)
		relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(relPath)
	}
}

func LsUntrackedHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	untrackedFiles, err := cfg.GetUntrackedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(d.LightBlue(d.Title("untracked files")))
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, uf := range untrackedFiles {
		absPath := filepath.Join(baseDir, uf)
		relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(relPath)
	}
}

func LsUntrackedTty(cfg *dotcfg.File) {
	untrackedFiles, err := cfg.GetUntrackedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("untracked files")
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, uf := range untrackedFiles {
		absPath := filepath.Join(baseDir, uf)
		relPath, err := dotcfg.GetRelativeToBaseDir(cwd, absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(relPath)
	}
}
