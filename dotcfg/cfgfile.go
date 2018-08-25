package dotcfg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

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

func LoadCfg(baseDir string) (*File, error) {
	filePath := filepath.Join(baseDir, FileName)

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open %v", filePath)
	}

	cfg := File{}
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse %v\n", filePath)
	}
	return &cfg, nil
}

func MustLoadCfg(baseDir string) *File {
	cfgFile, err := LoadCfg(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return cfgFile
}

func (f *File) MustSerialize(baseDir string, tx *rollback.Tx) {
	if err := f.Serialize(baseDir, tx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Serialize serializes the cfgfile at the given destination.
// If the destination is an empty string, cwd is used instead.
func (cfgFile *File) Serialize(baseDir string, tx *rollback.Tx) error {
	filePath := filepath.Join(baseDir, FileName)

	isCreated := false
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			isCreated = true
		} else {
			err = fmt.Errorf("failed to stat %v\n%v", filePath, err)
			txErr := tx.Rollback()
			if txErr != nil {
				return fmt.Errorf("during error:\n%v\ntransaction rollback failed with error:\n%v", err, txErr)
			}
			return err
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("failed to create or open %v\n%v", filePath, err)
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("during error:\n%v\ntransaction rollback failed with error:\n%v", err, txErr)
		}
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(cfgFile); err != nil {
		err = fmt.Errorf("failed to write %v\n%v", filePath, err)
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("during error:\n%v\ntransaction rollback failed with error:\n%v", err, txErr)
		}
		return err
	}

	// explicitly close f and check the error to ensure the file has been written
	if err := f.Close(); err != nil {
		err = fmt.Errorf("failed to close file %v\n%v", filePath, err)
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("during error:\n%v\ntransaction rollback failed with error:\n%v", err, txErr)
		}
		return err
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, filePath)
	}
	return nil
}

func mustStage(baseDir, filePath string) {
	if err := stage(baseDir, filePath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func stage(baseDir, filePath string) error {
	cmd := exec.Command("git", "-C", baseDir, "add", filePath)
	if err := cmdrunner.PipeTo(cmd, nil, os.Stderr); err != nil {
		return fmt.Errorf("failed to stage %v\n%v", filePath, err)
	}
	return nil
}

func mustCommit(baseDir, msg string) {
	if err := commit(baseDir, msg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func commit(baseDir, msg string) error {
	cmd := exec.Command("git", "-C", baseDir, "commit", "-m", msg)
	if err := cmdrunner.PipeTo(cmd, nil, os.Stderr); err != nil {
		cmdString := fmt.Sprintf("%v \"%v\"", strings.Join(cmd.Args[:3], " "), msg)
		return fmt.Errorf("'%v' failed with error:\n%v\n", cmdString, err)
	}
	return nil
}

func (cfg *File) MustStage(baseDir string) {
	if err := cfg.Stage(baseDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *File) Stage(baseDir string) error {
	for _, t := range cfg.Templates {
		if err := stage(baseDir, t.FilePath); err != nil {
			return err
		}
	}
	for _, i := range cfg.Instances {
		if err := stage(baseDir, i.FilePath); err != nil {
			return err
		}
	}
	for _, s := range cfg.Singletons {
		if err := stage(baseDir, s.FilePath); err != nil {
			return err
		}
	}
	schemasDirPath := filepath.Join(baseDir, SchemasDirName)
	if _, err := os.Stat(schemasDirPath); err != nil && os.IsNotExist(err) {
		// create .cfg_schemas if not exist
		// it can be incidentally rm'd by git, if it becomes an empty directory
		if err := os.MkdirAll(schemasDirPath, os.ModePerm); err != nil {
			return err
		}
	}
	if err := stage(baseDir, schemasDirPath); err != nil {
		return err
	}
	return cfg.StageSelf(baseDir)
}

func (cfg *File) Commit(baseDir string) error {
	// TODO: figure out changes and make appropriate message
	return commit(baseDir, "add changes")
}

func (cfg *File) MustCommit(baseDir string) {
	if err := cfg.Commit(baseDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *File) StageSelf(baseDir string) error {
	return stage(baseDir, ".cfg.json")
}

func (cfg *File) MustStageSelf(baseDir string) {
	if err := cfg.StageSelf(baseDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *File) CommitSelf(baseDir string) error {
	return commit(baseDir, "add .cfg.json")
}

type fileInfoTup struct {
	baseDir string
	fInfo   os.FileInfo
}

func (cfg *File) GetUntrackedFiles() ([]string, error) {
	trackedFiles := make(map[string]bool)
	for _, t := range cfg.Templates {
		trackedFiles[t.FilePath] = true
	}
	for _, i := range cfg.Instances {
		trackedFiles[i.FilePath] = true
	}
	for _, s := range cfg.Singletons {
		trackedFiles[s.FilePath] = true
	}

	// TODO: find cfg base dir
	// there are lots of other places in the code like this
	// which aren't marked with TODO comments
	baseDir := "."
	wdFileInfo, err := os.Stat(baseDir)
	if err != nil {
		return nil, err
	}

	untracked := make([]string, 0)
	stack := make([]fileInfoTup, 1)
	stack[0] = fileInfoTup{baseDir, wdFileInfo}
	for len(stack) > 0 {
		var s fileInfoTup
		s, stack = stack[len(stack)-1], stack[:len(stack)-1]
		filePath := filepath.Join(s.baseDir, s.fInfo.Name())

		// TODO: clean up
		// this could have bugs based on working directory
		// and relative paths
		if strings.HasPrefix(filePath, ".git") {
			continue
		}
		if strings.HasPrefix(filePath, ".cfg/") || strings.HasPrefix(filePath, ".cfg\\") {
			continue
		}
		if strings.HasPrefix(filePath, ".cfg_schemas/") || strings.HasPrefix(filePath, ".cfg_schemas\\") {
			continue
		}
		if filePath == ".cfg.json" {
			continue
		}

		if s.fInfo.IsDir() {
			fs, err := ioutil.ReadDir(filePath)
			if err != nil {
				return nil, err
			}
			for _, childF := range fs {
				stack = append(stack, fileInfoTup{filePath, childF})
			}
		} else {
			_, isTracked := trackedFiles[filePath]
			if !isTracked {
				untracked = append(untracked, filePath)
			}
		}
	}
	return untracked, nil
}

func (cfg *File) Infer(filePath string) error {
	// TODO: filePath is assumed to be a relative path right now

	// TODO: find cfg base dir
	// there are lots of other places in the code like this
	// which aren't marked with TODO comments

	// TODO: should we leave this as `exec` or should we import code???
	// leaving as `exec` means no infer features if shipping a single static binary
	// (need to have the `schema` binary as well)
	// on the other hand, it is a good way to increase `schema` usage
	dest := filepath.Join(SchemasDirName, filePath)

	parentDir := filepath.Dir(dest)
	if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpDest := fmt.Sprintf("%v.%v", dest, time.Now().UnixNano())

	cmd := exec.Command("schema", "infer", tmpDest, "--make-required")
	err = cmdrunner.PipeThrough(cmd, f, nil, nil)
	if err != nil {
		if rmErr := os.Remove(tmpDest); rmErr != nil {
			return fmt.Errorf("during this error:\n%v\nos.Remove failed:\n%v", err, rmErr)
		}
		return err
	}
	return os.Rename(tmpDest, dest)
}

func (cfg *File) MustRmSchema(target string, onlyFromIndex bool) {
	if err := cfg.RmSchema(target, onlyFromIndex); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *File) RmSchema(target string, onlyFromIndex bool) error {
	// TODO: target is assumed to be a relative path
	// TODO: assumed to be run from base directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	schemaPath := filepath.Join(cwd, SchemasDirName, target)

	if _, err := os.Stat(schemaPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if !cfg.NoGit {
		if onlyFromIndex {
			out, err := exec.Command("git", "rm", "--cached", schemaPath).CombinedOutput()
			if err != nil {
				return fmt.Errorf("'git rm --cached %v' failed\nerror: %v\noutput: %v", target, err, string(out))
			}
		} else {
			out, err := exec.Command("git", "rm", schemaPath).CombinedOutput()
			if err != nil {
				return fmt.Errorf("'git rm %v' failed\nerror: %v\noutput: %v", target, err, string(out))
			}
		}
	} else {
		if !onlyFromIndex {
			if err := os.Remove(schemaPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cfgFile *File) MustWarnDiffs(templName, instFilePath string, wErr io.Writer) {
	if err := cfgFile.WarnDiffs(templName, instFilePath, wErr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (cfgFile *File) WarnDiffs(templName, instFilePath string, wErr io.Writer) error {
	var templSchemaPath string
	for _, templ := range cfgFile.Templates {
		if templ.Name == templName {
			if templ.Schema.FilePath != "" {
				templSchemaPath = templ.Schema.FilePath
			} else {
				// TODO: relative dir problems
				templSchemaPath = filepath.Join(SchemasDirName, templ.FilePath)
			}
			break
		}
	}
	if templSchemaPath == "" {
		return fmt.Errorf("template '%v' does not exist", templName)
	}
	// TODO: relative dir problems
	instSchemaPath := filepath.Join(SchemasDirName, instFilePath)

	args := []string{
		"diff",
		templSchemaPath,
		instSchemaPath,
		"--miss-from-1",
		fmt.Sprintf("'%v' templ, yet is in '%v'", templName, instFilePath),
		"--miss-from-2",
		fmt.Sprintf("'%v', yet is in '%v' templ", instFilePath, templName),
		"--differ-1",
		fmt.Sprintf("the '%v' templ", templName),
		"--differ-2",
		fmt.Sprintf("'%v'", instFilePath),
	}
	cmd := exec.Command("schema", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		niceErr := fmt.Errorf("'%v' failed: %v", strings.Join(cmd.Args[:4], " "), err)
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return niceErr
		}
		status, ok := exiterr.Sys().(syscall.WaitStatus)
		if !ok {
			return niceErr
		}
		if status.ExitStatus() != 2 {
			return niceErr
		}
		fmt.Fprintf(wErr, "%v", string(out))
	}
	return nil
}
