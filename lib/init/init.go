package init

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"

	"github.com/confbase/cfg/lib/dotcfg"
)

func Init(append, overwrite bool) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}

	dotCfgDirPath := mkDotCfgDir(cwd)
	mkDotfile(cwd)
	mkKeyfile(dotCfgDirPath)
	mkGitIgnore(cwd, append, overwrite)

	fmt.Printf("Initialized empty base in %v\n", dotCfgDirPath)
}

func mkDotCfgDir(baseDir string) string {
	dirPath := path.Join(baseDir, dotcfg.Dirname)
	notExistErrOut(dirPath, "")

	if err := os.Mkdir(dirPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create directory %v\n", dirPath)
		os.Exit(1)
	}
	return dirPath
}

func mkDotfile(baseDir string) {
	filePath := path.Join(baseDir, dotcfg.FileName)
	notExistErrOut(filePath, "")

	dotfile := dotcfg.File{
		Templates: make([]string, 0),
		Instances: make(map[string]string),
	}
	data, err := json.Marshal(dotfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to serialize %v\n", filePath)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write to %v\n", filePath)
		os.Exit(1)
	}
}

func mkKeyfile(baseDir string) {
	filePath := path.Join(baseDir, dotcfg.KeyfileName)
	notExistErrOut(filePath, "")

	keyfile := dotcfg.Key{
		Email:      viper.GetString("email"),
		EntryPoint: viper.GetString("entryPoint"),
	}
	data, err := json.Marshal(keyfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to serialize %v\n", filePath)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write to %v\n", filePath)
		os.Exit(1)
	}
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
		notExistErrOut(filePath, "did you mean to use --overwrite-gitignore or --append-to-gitignore?\n")
	}

	if err := ioutil.WriteFile(filePath, []byte(ignoreStr), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write to %v\n", filePath)
		os.Exit(1)
	}
}

func notExistErrOut(filePath, msg string) {
	_, err := os.Stat(filePath)
	if err == nil || (err != nil && !os.IsNotExist(err)) {
		fmt.Fprintf(os.Stderr, "error: %v already exists\n", filePath)
		fmt.Fprintf(os.Stderr, "%v", msg)
		os.Exit(1)
	}
}
