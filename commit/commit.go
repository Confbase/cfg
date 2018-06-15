package commit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
)

func MustCommit(cfg *Config) {
	if err := Commit(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func gitAdd(filePath string) error {
	cmd := exec.Command("git", "add", filePath)
	return cmdrunner.PipeTo(cmd, os.Stdout, os.Stderr)
}

func Commit(cfg *Config) error {
	out, err := exec.Command("git", "diff", "--name-only", "--cached").CombinedOutput()
	if err != nil {
		return fmt.Errorf("'git diff --name-only --cached' failed with error:\n%vand output:\n%v", err, string(out))
	}
	commitCmd := exec.Command("git", "commit", "-m", cfg.Message)
	if len(out) != 0 {
		// there are already files in the staging area
		return cmdrunner.PipeTo(commitCmd, os.Stdout, os.Stderr)
	}

	// must add tracked files to staging, then commit
	cfgFile, err := dotcfg.LoadCfg()
	if err != nil {
		return err
	}

	for _, templ := range cfgFile.Templates {
		gitAdd(templ.FilePath)
		if templ.Schema.FilePath != "" {
			if err := gitAdd(templ.Schema.FilePath); err != nil {
				return err
			}
		} else {
			if err := gitAdd(filepath.Join(dotcfg.SchemasDirName, templ.FilePath)); err != nil {
				return err
			}
		}
	}

	for _, inst := range cfgFile.Instances {
		gitAdd(inst.FilePath)
		if inst.Schema.FilePath != "" {
			if err := gitAdd(inst.Schema.FilePath); err != nil {
				return err
			}
		} else {
			if err := gitAdd(filepath.Join(dotcfg.SchemasDirName, inst.FilePath)); err != nil {
				return err
			}
		}
	}

	for _, s := range cfgFile.Singletons {
		gitAdd(s.FilePath)
		// should we add singleton schemas?
		// seems the answer is yes, but what about singletons like .gitignore?
/*
		if s.Schema.FilePath != "" {
			if err := gitAdd(s.Schema.FilePath); err != nil {
				return err
			}
		} else {
			if err := gitAdd(filepath.Join(dotcfg.SchemasDirName, s.FilePath)); err != nil {
				return err
			}
		}
*/
	}

	return cmdrunner.PipeTo(commitCmd, os.Stdout, os.Stderr)
}
