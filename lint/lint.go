package lint

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/dotcfg"
)

func MustLint() {
	if err := Lint(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Lint() error {
	cfgFile, err := dotcfg.LoadCfg()
	if err != nil {
		return err
	}
	for _, templ := range cfgFile.Templates {
		for _, inst := range cfgFile.Instances {
			isInstOfTempl := false
			for _, templName := range inst.TemplNames {
				if templName == templ.Name {
					isInstOfTempl = true
					break
				}
			}
			if isInstOfTempl {
				if err := cfgFile.WarnDiffs(templ.Name, inst.FilePath, os.Stderr); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
