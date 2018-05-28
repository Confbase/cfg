package clone

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Confbase/cfg/cmdrunner"
)

func Clone(cloneCfg *Config) {
	if cloneCfg.NoGit {
		cloneNoGit(cloneCfg)
	} else {
		cloneWithGit(cloneCfg)
	}
}

func cloneNoGit(cloneCfg *Config) {

}

func cloneWithGit(cloneCfg *Config) {
	var targetDir string
	if cloneCfg.TargetDir != "" {
		targetDir = cloneCfg.TargetDir
	} else {
		elems := strings.Split(cloneCfg.RemoteURL, "/")
		if len(elems) <= 1 {
			elems = strings.Split(cloneCfg.RemoteURL, ":")
		}
		targetDir = elems[len(elems)-1]
	}

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chdir(targetDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmdrunner.RunOrFatal(exec.Command("cfg", "init"))
	cmdrunner.RunOrFatal(exec.Command("cfg", "remote", "add", "origin", cloneCfg.RemoteURL))
	if cloneCfg.DefaultSnap != "master" {
		cmdrunner.RunOrFatal(exec.Command("cfg", "snap", "new", cloneCfg.DefaultSnap))
	}
	cmdrunner.RunOrFatal(exec.Command("git", "fetch", "origin", cloneCfg.DefaultSnap))
	remoteBranch := fmt.Sprintf("origin/%v", cloneCfg.DefaultSnap)
	cmdrunner.RunOrFatal(exec.Command("git", "reset", "--hard", remoteBranch))
}
