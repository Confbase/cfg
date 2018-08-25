package snap

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
)

func Remove(targets []string) {
	snapsFile := dotcfg.MustLoadSnaps()
	targetSet := make(map[string]bool)
	for _, target := range targets {
		if target == snapsFile.Current.Name {
			fmt.Fprintf(os.Stderr, "error: cannot rm the snap which is checked out\n")
			os.Exit(1)
		}
		targetSet[target] = true
	}
	newSnaps := make([]dotcfg.Snapshot, 0)
	for _, s := range snapsFile.Snapshots {
		if _, ok := targetSet[s.Name]; !ok {
			newSnaps = append(newSnaps, s)
		}
	}

	cfgFile := dotcfg.MustLoadCfg("")
	if !cfgFile.NoGit {
		for target, _ := range targetSet {
			cmdrunner.RunOrFatal(exec.Command("git", "branch", "-D", target))
		}
		cfgFile.MustSerialize("", nil)
	}

	snapsFile.Snapshots = newSnaps
	snapsFile.MustSerialize("", nil)
}
