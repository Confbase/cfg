package init

import (
	"fmt"
	"io/ioutil"
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
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory\n%v", err)
	}

	filePath := filepath.Join(cwd, dotcfg.FileName)
	dirPath := filepath.Join(cwd, dotcfg.DirName)

	if err := existsErrOut(filePath, "", nil); err != nil {
		return err
	}
	if err := existsErrOut(dirPath, "", nil); err != nil {
		return err
	}

	tx := rollback.NewTx()

	if !cfg.NoModGitIgnore {
		if err := mkGitIgnore(cwd, cfg, tx); err != nil {
			return err
		}
	}

	cfgFile := dotcfg.NewCfg()
	cfgFile.NoGit = cfg.NoGit
	if err := cfgFile.Serialize(tx); err != nil {
		return err
	}

	keyfile := dotcfg.NewKey(filepath.Base(cwd))
	if err := keyfile.Serialize(tx); err != nil {
		return err
	}

	snaps := dotcfg.NewSnaps()
	if err := snaps.Serialize(tx); err != nil {
		return err
	}

	if err := mkSchemasDir(cwd, tx); err != nil {
		return err
	}

	if !cfgFile.NoGit {
		if err := initGitRepo(cwd, cfg.Force, tx); err != nil {
			return err
		}
		if err := cfgFile.StageSelf(); err != nil {
			return err
		}
		if err := cfgFile.CommitSelf(); err != nil {
			return err
		}
		// TODO: check err on track.Track
		track.Track(".gitignore")
	}

	fmt.Printf("Initialized empty base in %v\n", cwd)
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
		existsCmd := exec.Command("git", "rev-parse", "--git-dir")
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

	cmd := exec.Command("git", "init")
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
			err = fmt.Errorf("failed to open .gitignore for appending\n%v", err)
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
			return err
		}
		defer f.Close()

		if _, err = f.WriteString(ignoreStr); err != nil {
			err = fmt.Errorf("failed to write to .gitignore\n")
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

	if err := ioutil.WriteFile(filePath, []byte(ignoreStr), 0644); err != nil {
		err = fmt.Errorf("failed to write to %v\n%v", filePath, err)
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
