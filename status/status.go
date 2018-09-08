package status

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
)

func MustShowStatus() {
	if err := ShowStatus(os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func ShowStatus(w, wErr io.Writer) error {
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		return err
	}
	cfgFile, err := dotcfg.LoadCfg(baseDir)
	if err != nil {
		return err
	}
	if !cfgFile.NoGit {
		return cmdrunner.PipeTo(exec.Command("git", "status", "-sb"), w, wErr)
	}
	// TODO: --no-git
	return nil
}
