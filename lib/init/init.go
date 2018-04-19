package init

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/Confbase/cfg/lib/dotcfg"
	"github.com/Confbase/cfg/lib/rollback"
	"github.com/Confbase/cfg/lib/util"
)

func Init(appendGitIgnore, overwriteGitIgnore, noGit, noModGitIgnore bool) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}

	filePath := path.Join(cwd, dotcfg.FileName)
	dirPath := path.Join(cwd, dotcfg.Dirname)

	existsErrOut(filePath, "", nil)
	existsErrOut(dirPath, "", nil)

	tx := rollback.NewTx()

	if !noModGitIgnore {
		mkGitIgnore(cwd, appendGitIgnore, overwriteGitIgnore, tx)
	}

	cfg := dotcfg.NewCfg()
	cfg.NoGit = noGit
	cfg.MustSerialize(tx)

	keyfile := dotcfg.NewKey()
	keyfile.MustSerialize(tx)

	snaps := dotcfg.NewSnaps()
	snaps.MustSerialize(tx)

	if !cfg.NoGit {
		initGitRepo(cwd, tx)
		cfg.MustStageSelf()
		cfg.MustCommitSelf()
		cfg.MustAddGitIgnore()
	}

	fmt.Printf("Initialized empty base in %v\n", cwd)
}

func initGitRepo(baseDir string, tx *rollback.Tx) {
	existsCmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := existsCmd.Run(); err == nil {
		fmt.Fprintf(os.Stderr, "error: the current directory is inside a git repository\n")
		fmt.Fprintf(os.Stderr, "if you wish to use cfg in conjuction with an existing git repository,\n")
		fmt.Fprintf(os.Stderr, "consider running 'cfg init --no-git'\n")
		tx.MustRollback()
		os.Exit(1)
	}

	cmd := exec.Command("git", "init")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		tx.MustRollback()
		os.Exit(1)
	}

	tx.DirsCreated = append(tx.DirsCreated, path.Join(baseDir, ".git"))
}

func mkGitIgnore(baseDir string, appendGitIgnore, overwriteGitIgnore bool, tx *rollback.Tx) {
	filePath := path.Join(baseDir, ".gitignore")
	ignoreStr := ".cfg/\n"

	if appendGitIgnore {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open .gitignore for appending\n")
			tx.MustRollback()
			os.Exit(1)
		}
		defer f.Close()

		if _, err = f.WriteString(ignoreStr); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to .gitignore\n")
			tx.MustRollback()
			os.Exit(1)
		}
		return
	}

	if !overwriteGitIgnore {
		existsErrOut(filePath, "you must specify --overwrite-gitignore, --append-to-gitignore, or --no-modify-gitignore\n", tx)
	}

	isCreated := false
	isExist, err := util.Exists(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to stat %v\n", filePath)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		tx.MustRollback()
		os.Exit(1)
	}
	if !isExist {
		isCreated = true
	}

	if err := ioutil.WriteFile(filePath, []byte(ignoreStr), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write to %v\n", filePath)
		tx.MustRollback()
		os.Exit(1)
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, filePath)
	}
}

func existsErrOut(filePath, msg string, tx *rollback.Tx) {
	_, err := os.Stat(filePath)
	if err == nil || (err != nil && !os.IsNotExist(err)) {
		fmt.Fprintf(os.Stderr, "error: %v already exists\n", filePath)
		fmt.Fprintf(os.Stderr, "%v", msg)
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}
}
