package init

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Confbase/cfg/dotcfg"
	"github.com/Confbase/cfg/rollback"
	"github.com/Confbase/cfg/track"
)

func MustInit(cfg *Config) {
	if err := Init(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Init(cfg *Config) error {
	var dest string

	if len(cfg.Dest) == 0 || cfg.Dest[0] != '/' {
		// if the given destination is empty or a relative path,
		// then prepend with cwd
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory\n%v", err)
		}
		dest = filepath.Join(cwd, cfg.Dest)
	} else {
		// otherwise, interpret the destination as an absolute path
		dest = cfg.Dest
	}

	if _, err := os.Stat(dest); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// if the destination doesn't exist, make the directories
		if err := os.MkdirAll(dest, 0755); err != nil {
			return err
		}
	}

	filePath := filepath.Join(dest, dotcfg.FileName)
	dirPath := filepath.Join(dest, dotcfg.DirName)

	if err := existsErrOut(filePath, "", nil); err != nil {
		return err
	}
	if err := existsErrOut(dirPath, "", nil); err != nil {
		return err
	}

	tx := rollback.NewTx()

	if !cfg.NoModGitIgnore {
		if err := mkGitIgnore(dest, cfg, tx); err != nil {
			return err
		}
	}

	cfgFile := dotcfg.NewCfg()
	cfgFile.NoGit = cfg.NoGit
	if err := cfgFile.Serialize(dest, tx); err != nil {
		return err
	}

	// use the basename of the destination as the name of this base
	keyfile := dotcfg.NewKey(filepath.Base(dest))
	if err := keyfile.Serialize(dest, tx); err != nil {
		return err
	}

	snaps := dotcfg.NewSnaps()
	if err := snaps.Serialize(dest, tx); err != nil {
		return err
	}

	if err := mkSchemasDir(dest, tx); err != nil {
		return err
	}

	if !cfgFile.NoGit {
		if err := initGitRepo(dest, cfg.Force, tx); err != nil {
			return err
		}
		if err := cfgFile.StageSelf(dest); err != nil {
			return err
		}
		if err := cfgFile.CommitSelf(dest); err != nil {
			return err
		}
		// TODO: check err on track.Track
		track.Track(dest, filepath.Join(dest, ".gitignore"))
	}

	fmt.Printf("Initialized empty base at %v\n", dest)
	return nil
}

func mkSchemasDir(baseDir string, tx *rollback.Tx) error {
	schemasDir := filepath.Join(baseDir, dotcfg.SchemasDirName)
	if err := os.MkdirAll(schemasDir, os.ModePerm); err != nil {
		return err
	}

	tx.DirsCreated = append(tx.DirsCreated, schemasDir)
	return nil
}

func initGitRepo(baseDir string, force bool, tx *rollback.Tx) error {
	if !force {
		existsCmd := exec.Command("git", "-C", baseDir, "rev-parse", "--git-dir")
		if err := existsCmd.Run(); err == nil {
			err = fmt.Errorf("the current directory is inside a git repository\n")
			err = fmt.Errorf("%v if you wish to use cfg in conjuction with an existing git repository,\n", err)
			err = fmt.Errorf("%v consider running 'cfg init --no-git'\n", err)
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
			return err
		}
	}

	cmd := exec.Command("git", "init", baseDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("'git init' failed\nerror: %v\noutput: %v\n", err, string(out))
		if txErr := tx.Rollback(); txErr != nil {
			err = rollback.MergeTxErr(err, txErr)
		}
		return err
	}

	tx.DirsCreated = append(tx.DirsCreated, filepath.Join(baseDir, ".git"))
	return nil
}

func mkGitIgnore(baseDir string, cfg *Config, tx *rollback.Tx) error {
	filePath := filepath.Join(baseDir, ".gitignore")
	ignoreStr := ".cfg/\n"

	if cfg.AppendGitIgnore {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			err = fmt.Errorf("failed to open %v for appending\n%v", filePath, err)
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
			return err
		}
		defer f.Close()

		if _, err = f.WriteString(ignoreStr); err != nil {
			err = fmt.Errorf("failed to write to %v\n", filePath)
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
			return err
		}

		// explicitly close f to check the error
		if err = f.Close(); err != nil {
			err = fmt.Errorf("failed to close to %v\n", filePath)
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
			return err
		}
		return nil
	}

	if !cfg.OverwriteGitIgnore {
		msg := "you must specify --overwrite-gitignore,"
		msg = fmt.Sprintf("%v --append-to-gitignore,", msg)
		msg = fmt.Sprintf("%v or --no-modify-gitignore\n", msg)
		if err := existsErrOut(filePath, msg, tx); err != nil {
			return err
		}
	}

	isCreated := false
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			isCreated = true
		} else {
			err = fmt.Errorf("failed to stat %v\n%v", filePath, err)
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
			return err
		}
	}

	// TODO: refactor following three function calls to adhere to DRY
	// (exact same code lives one page above)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		err = fmt.Errorf("failed to open %v\n%v", filePath, err)
		if txErr := tx.Rollback(); txErr != nil {
			err = rollback.MergeTxErr(err, txErr)
		}
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(ignoreStr); err != nil {
		err = fmt.Errorf("failed to write to %v\n", filePath)
		if txErr := tx.Rollback(); txErr != nil {
			err = rollback.MergeTxErr(err, txErr)
		}
		return err
	}

	// explicitly close f to check the error
	if err = f.Close(); err != nil {
		err = fmt.Errorf("failed to close to %v\n", filePath)
		if txErr := tx.Rollback(); txErr != nil {
			err = rollback.MergeTxErr(err, txErr)
		}
		return err
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, filePath)
	}
	return nil
}

func existsErrOut(filePath, msg string, tx *rollback.Tx) error {
	_, err := os.Stat(filePath)
	if err == nil {
		err = fmt.Errorf("%v already exists\n%v", filePath, msg)
	} else if err != nil && !os.IsNotExist(err) {
		err = fmt.Errorf("failed to stat %v\n%v", filePath, err)
	} else {
		err = nil
	}

	if err != nil && tx != nil {
		if txErr := tx.Rollback(); txErr != nil {
			err = rollback.MergeTxErr(err, txErr)
		}
	}
	return err
}
