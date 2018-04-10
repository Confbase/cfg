package ls

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/confbase/cfg/lib/decorate"
	"github.com/confbase/cfg/lib/dotcfg"
)

func isStdoutTty() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func Ls() {
	cfg := dotcfg.MustLoadCfg()

	if isStdoutTty() {
		LsTemplHuman(cfg)
		LsInstancesHuman(cfg)
		LsSingletonsHuman(cfg)
	} else {
		LsTemplTty(cfg)
		LsInstancesTty(cfg)
		LsSingletonsTty(cfg)
	}
}

func LsTemplHuman(cfg *dotcfg.File) {
	fmt.Println(decorate.LightBlue(decorate.Title("templates")))
	if len(cfg.Templates) > 0 {
		for _, t := range cfg.Templates {
			fmt.Printf(decorate.Green("%v")+": %v\n", t.Name, t.FilePath)
		}
	}
	fmt.Println()
}

func LsTemplTty(cfg *dotcfg.File) {
	fmt.Println("templates")
	for _, t := range cfg.Templates {
		fmt.Printf("%v\t%v\t%v\n", t.Name, t.FileType, t.FilePath)
	}
}

func LsInstancesHuman(cfg *dotcfg.File) {
	fmt.Println(decorate.LightBlue(decorate.Title("instances")))
	if len(cfg.Templates) > 0 {
		for templ, instances := range cfg.Instances {
			for _, i := range instances {
				fmt.Printf(decorate.Green("%v")+": %v\n", i, templ)
			}
		}
	}
	fmt.Println()
}

func LsInstancesTty(cfg *dotcfg.File) {
	fmt.Println("instances")
	for templ, instances := range cfg.Instances {
		for _, i := range instances {
			fmt.Printf("%v\t%v\n", templ, i)
		}
	}
}

func LsSingletonsHuman(cfg *dotcfg.File) {
	fmt.Println(decorate.LightBlue(decorate.Title("singletons")))
	for _, s := range cfg.Singletons {
		fmt.Println(s)
	}
}

func LsSingletonsTty(cfg *dotcfg.File) {
	fmt.Println("singletons")
	for _, s := range cfg.Singletons {
		fmt.Println(s)
	}
}
