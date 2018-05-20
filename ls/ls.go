package ls

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Confbase/cfg/decorate"
	"github.com/Confbase/cfg/dotcfg"
)

func isStdoutTty() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func Ls(noTty, noColors bool) {
	cfg := dotcfg.MustLoadCfg()

	if isStdoutTty() && !noTty {
		d := decorate.New()
		d.Enabled = !noColors

		snaps := dotcfg.MustLoadSnaps()
		fmt.Printf("## %v\n", snaps.Current)

		LsTemplHuman(cfg, d)
		LsInstancesHuman(cfg, d)
		LsSingletonsHuman(cfg, d)
	} else {
		LsTemplTty(cfg)
		LsInstancesTty(cfg)
		LsSingletonsTty(cfg)
	}
}

func LsTemplHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("templates")))
	if len(cfg.Templates) > 0 {
		for _, t := range cfg.Templates {
			fmt.Printf(d.Green("%v")+": %v\n", t.Name, t.FilePath)
		}
	}
	fmt.Println()
}

func LsTemplTty(cfg *dotcfg.File) {
	fmt.Println("templates")
	for _, t := range cfg.Templates {
		fmt.Printf("%v\t%v\n", t.Name, t.FilePath)
	}
}

func LsInstancesHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("instances")))
	if len(cfg.Templates) > 0 {
		for templ, instances := range cfg.Instances {
			for _, i := range instances {
				fmt.Printf(d.Green("%v")+": %v\n", i.FilePath, templ)
			}
		}
	}
	fmt.Println()
}

func LsInstancesTty(cfg *dotcfg.File) {
	fmt.Println("instances")
	for templ, instances := range cfg.Instances {
		for _, i := range instances {
			fmt.Printf("%v\t%v\n", templ, i.FilePath)
		}
	}
}

func LsSingletonsHuman(cfg *dotcfg.File, d *decorate.Decorator) {
	fmt.Println(d.LightBlue(d.Title("singletons")))
	for _, s := range cfg.Singletons {
		fmt.Println(s.FilePath)
	}
}

func LsSingletonsTty(cfg *dotcfg.File) {
	fmt.Println("singletons")
	for _, s := range cfg.Singletons {
		fmt.Println(s.FilePath)
	}
}
