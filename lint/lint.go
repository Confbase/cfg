package lint

import "github.com/Confbase/cfg/dotcfg"

func Lint() {
	cfgFile := dotcfg.MustLoadCfg()
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
				cfgFile.MustWarnDiffs(templ.Name, inst.FilePath)
			}
		}
	}
}
