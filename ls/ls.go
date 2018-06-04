package ls

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Confbase/cfg/decorate"
	"github.com/Confbase/cfg/dotcfg"
)

func isStdoutTty() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func Ls(lsCfg *Config) {
	cfg := dotcfg.MustLoadCfg()

	if isStdoutTty() && !lsCfg.NoTty {
		d := decorate.New()
		d.Enabled = !lsCfg.NoColors

		snaps := dotcfg.MustLoadSnaps()
		fmt.Printf("## %v\n", snaps.Current.Name)

		if !(lsCfg.DoLsTempls || lsCfg.DoLsInsts || lsCfg.DoLsSingles) {
			LsTemplsHuman(cfg, d)
			fmt.Println()
			LsInstsHuman(cfg, d)
			fmt.Println()
			LsSinglesHuman(cfg, d)
			return
		}
		if lsCfg.DoLsTempls {
			LsTemplsHuman(cfg, d)
			if lsCfg.DoLsInsts || lsCfg.DoLsSingles {
				fmt.Println()
			}
		}
		if lsCfg.DoLsInsts {
			LsInstsHuman(cfg, d)
			if lsCfg.DoLsSingles {
				fmt.Println()
			}
		}
		if lsCfg.DoLsSingles {
			LsSinglesHuman(cfg, d)
		}
	} else {
		if !(lsCfg.DoLsTempls || lsCfg.DoLsInsts || lsCfg.DoLsSingles) {
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
	}
}

func LsTemplsHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("templates")))
	if len(cfg.Templates) > 0 {
		for i, t := range cfg.Templates {
			end := "\n"
			if i == len(cfg.Templates)-1 {
				end = ""
			}
			fmt.Printf(d.Green("%v")+": %v%v", t.Name, t.FilePath, end)
		}
	}
	fmt.Println()
}

func LsTemplsTty(cfg *dotcfg.File) {
	fmt.Println("templates")
	for _, t := range cfg.Templates {
		fmt.Printf("%v\t%v\n", t.Name, t.FilePath)
	}
}

func LsInstsHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("instances")))
	if len(cfg.Templates) > 0 {
		for i, inst := range cfg.Instances {
			end := "\n"
			if i == len(cfg.Instances)-1 {
				end = ""
			}
			templsStr := strings.Join(inst.TemplNames, ", ")
			fmt.Printf(d.Green("%v")+": %v%v", inst.FilePath, templsStr, end)
		}
	}
	fmt.Println()
}

func LsInstsTty(cfg *dotcfg.File) {
	fmt.Println("instances")
	for _, inst := range cfg.Instances {
		templsStr := strings.Join(inst.TemplNames, ",")
		fmt.Printf("%v\t%v\n", inst.FilePath, templsStr)
	}
}

func LsSinglesHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("singletons")))
	for _, s := range cfg.Singletons {
		fmt.Println(s.FilePath)
	}
}

func LsSinglesTty(cfg *dotcfg.File) {
	fmt.Println("singletons")
	for _, s := range cfg.Singletons {
		fmt.Println(s.FilePath)
	}
}
