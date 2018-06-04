package dotcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Confbase/cfg/cmdrunner"
	"github.com/Confbase/cfg/rollback"
)

func NewInstance(filePath string) *Instance {
	return &Instance{
		FilePath:   filePath,
		TemplNames: make([]string, 0),
	}
}

func NewCfg() *File {
	return &File{
		Templates:  make([]Template, 0),
		Instances:  make([]Instance, 0),
		Singletons: make([]Singleton, 0),
		NoGit:      false,
	}
}

func MustLoadCfg() *File {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	filePath := filepath.Join(cwd, FileName)

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to open %v\n", filePath)
		os.Exit(1)
	}

	cfg := File{}
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse %v\n", filePath)
		os.Exit(1)
	}
	return &cfg
}

func (f *File) MustSerialize(tx *rollback.Tx) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}
	filePath := filepath.Join(cwd, FileName)

	cfgBytes, err := json.Marshal(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal key\n")
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}

	isCreated := false
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			isCreated = true
		} else {
			fmt.Fprintf(os.Stderr, "error: failed to stat %v\n", filePath)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			if tx != nil {
				tx.MustRollback()
			}
			os.Exit(1)
		}
	}

	if err := ioutil.WriteFile(filePath, cfgBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write %v\n", filePath)
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, filePath)
	}
}

func mustStage(filePath string) {
	cmd := exec.Command("git", "add", filePath)
	if err := cmdrunner.PipeFromCmd(cmd, nil, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to stage %v\n", filePath)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func mustCommit(msg string) {
	cmd := exec.Command("git", "commit", "-m", msg)
	if err := cmdrunner.PipeFromCmd(cmd, nil, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *File) MustStage() {
	for _, t := range cfg.Templates {
		mustStage(t.FilePath)
	}
	for _, i := range cfg.Instances {
		mustStage(i.FilePath)
	}
	for _, s := range cfg.Singletons {
		mustStage(s.FilePath)
	}
	cfg.MustStageSelf()
}

func (cfg *File) MustCommit() {
	// TODO: figure out changes and make appropriate message
	msg := "add changes"
	mustCommit(msg)
}

func (cfg *File) MustStageSelf() {
	mustStage(".cfg.json")
}

func (cfg *File) MustCommitSelf() {
	mustCommit("add .cfg.json")
}
