package commit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/dotcfg"
)

func MustCommit(cfg *Config) {
	if err := Commit(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func MustCommitOrRevert(cfg *Config) {
	if err := Commit(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		fmt.Fprintf(os.Stderr, "--- running 'git reset'\n")
		if err := cmdrunner.PipeTo(exec.Command("git", "reset"), os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}

func gitAdd(filePath string) error {
	cmd := exec.Command("git", "add", filePath)
	return cmdrunner.PipeTo(cmd, os.Stdout, os.Stderr)
}

// Commit could modify the staging area before returning an error.
// Callers should be careful to revert the staging area, if necessary.
func Commit(cfg *Config) error {
	out, err := exec.Command("git", "diff", "--name-only", "--cached").CombinedOutput()
	if err != nil {
		return fmt.Errorf("'git diff --name-only --cached' failed with error:\n%vand output:\n%v", err, string(out))
	}
	commitCmd := exec.Command("git", "commit", "-m", cfg.Message)
	if len(out) != 0 {
		// there are already files in the staging area
		// TODO: infer schemas and before committing
		return cmdrunner.PipeTo(commitCmd, os.Stdout, os.Stderr)
	}

	// must add tracked files to staging, then commit
	cfgFile, err := dotcfg.LoadCfg()
	if err != nil {
		return err
	}

	for _, templ := range cfgFile.Templates {
		if err := gitAdd(templ.FilePath); err != nil {
			return err
		}
		if templ.Schema.FilePath != "" {
			if err := gitAdd(templ.Schema.FilePath); err != nil {
				return err
			}
		} else {
			// TODO: add config variable + data structure to only do this to certain files
			if err := cfgFile.Infer(templ.FilePath); err != nil {
				exiterr, ok := err.(*exec.ExitError)
				if !ok {
					return err
				}
				// 'schema diff' exited with an exit code != 0
				status, ok := exiterr.Sys().(syscall.WaitStatus)
				if !ok {
					return err
				}
				if status.ExitStatus() == 1 || status.ExitStatus() == 127 {
					continue
				}
				return err
			}
			if err := gitAdd(filepath.Join(dotcfg.SchemasDirName, templ.FilePath)); err != nil {
				return err
			}
		}
	}

	for _, inst := range cfgFile.Instances {
		if err := gitAdd(inst.FilePath); err != nil {
			return err
		}
		if inst.Schema.FilePath != "" {
			if err := gitAdd(inst.Schema.FilePath); err != nil {
				return err
			}
		} else {
			// TODO: add config variable + data structure to only do this to certain files
			if err := cfgFile.Infer(inst.FilePath); err != nil {
				return err
			}
			for _, templName := range inst.TemplNames {
				for _, templ := range cfgFile.Templates {
					if templ.Name == templName {
						var errBuff bytes.Buffer
						if cfgFile.WarnDiffs(templName, inst.FilePath, &errBuff); err != nil {
							exiterr, ok := err.(*exec.ExitError)
							if !ok {
								return err
							}
							// 'schema diff' exited with an exit code != 0
							status, ok := exiterr.Sys().(syscall.WaitStatus)
							if !ok {
								return err
							}
							if status.ExitStatus() == 1 || status.ExitStatus() == 127 {
								continue
							}
							return err
						}
						if errBuff.Len() != 0 {
							return fmt.Errorf("--- 'cfg lint' gave output:\n%v", errBuff.String())
						}
					}
				}
			}
			if err := gitAdd(filepath.Join(dotcfg.SchemasDirName, inst.FilePath)); err != nil {
				return err
			}
		}
	}

	for _, s := range cfgFile.Singletons {
		if err := gitAdd(s.FilePath); err != nil {
			return err
		}
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
