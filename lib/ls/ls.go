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
	templStr := decorate.LightBlue(decorate.Title("templates\n"))
	if len(cfg.Templates) > 0 {
		for _, t := range cfg.Templates {
			templStr = fmt.Sprintf(
				"%v"+decorate.Green("%v: %v\n"),
				templStr,
				t.Name,
				t.FilePath,
			)
		}
	}
	fmt.Printf("%v\n", templStr)
}

func LsTemplTty(cfg *dotcfg.File) {
	fmt.Println("templates")
	for _, t := range cfg.Templates {
		fmt.Printf("%v\t%v\t%v\n", t.Name, t.FileType, t.FilePath)
	}
}

func LsInstancesHuman(cfg *dotcfg.File) {
	instancesStr := decorate.LightBlue(decorate.Title("instances\n"))
	if len(cfg.Templates) > 0 {
		for templ, instances := range cfg.Instances {
			for _, i := range instances {
				instancesStr = fmt.Sprintf(
					"%v"+decorate.Green("%v: %v\n"),
					instancesStr,
					i,
					templ,
				)
			}
		}
	}
	fmt.Printf("%v\n", instancesStr)
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
	singletonsStr := decorate.LightBlue(decorate.Title("singletons"))
	fmt.Printf("%v\n", singletonsStr)
}

func LsSingletonsTty(cfg *dotcfg.File) {
	fmt.Printf("singletons")
}
