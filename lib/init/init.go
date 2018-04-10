package init

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"

	"github.com/Confbase/cfg/lib/dotcfg"
)

func Init(append, overwrite bool) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}

	filePath := path.Join(cwd, dotcfg.FileName)
	dirPath := path.Join(cwd, dotcfg.Dirname)
	keyPath := path.Join(dirPath, dotcfg.KeyfileName)
	existsErrOut(filePath, "")
	existsErrOut(dirPath, "")
	existsErrOut(keyPath, "")

	dotfile := dotcfg.File{
		Templates: make([]dotcfg.Template, 0),
		Instances: make(map[string]([]string)),
	}
	dotfile.MustSerialize()

	keyfile := dotcfg.Key{
		Remotes:    make(map[string]string),
		Email:      viper.GetString("email"),
		EntryPoint: viper.GetString("entryPoint"),
	}
	keyfile.MustSerialize()

	mkGitIgnore(cwd, append, overwrite)

	fmt.Printf("Initialized empty base in %v\n", cwd)
}

func mkGitIgnore(baseDir string, append, overwrite bool) {
	filePath := path.Join(baseDir, ".gitignore")
	ignoreStr := ".cfg/\n"

	if append {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open .gitignore for appending\n")
			os.Exit(1)
		}
		defer f.Close()

		if _, err = f.WriteString(ignoreStr); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to .gitignore\n")
			os.Exit(1)
		}
		return
	}

	if !overwrite {
		existsErrOut(filePath, "did you mean to use --overwrite-gitignore or --append-to-gitignore?\n")
	}

	if err := ioutil.WriteFile(filePath, []byte(ignoreStr), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write to %v\n", filePath)
		os.Exit(1)
	}
}

func existsErrOut(filePath, msg string) {
	_, err := os.Stat(filePath)
	if err == nil || (err != nil && !os.IsNotExist(err)) {
		fmt.Fprintf(os.Stderr, "error: %v already exists\n", filePath)
		fmt.Fprintf(os.Stderr, "%v", msg)
		os.Exit(1)
	}
}
