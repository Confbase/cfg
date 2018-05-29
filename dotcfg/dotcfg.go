package dotcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/viper"

	"github.com/Confbase/cfg/rollback"
)

const (
	Dirname       = ".cfg"       // this dir resides in ./
	FileName      = ".cfg.json"  // this file resides in ./
	KeyfileName   = "key.json"   // this file resides in ./.cfg/
	SnapsFileName = "snaps.json" // this file resides in ./.cfg/
)

type Template struct {
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
}

// .cfg.json is tracked by git
type File struct {
	Templates  []Template              `json:"templates"`
	Instances  map[string]([]Instance) `json:"instances"`
	Singletons []Singleton             `json:"singletons"`
	NoGit      bool                    `json:"noGit"`
}

type Singleton struct {
	FilePath string `json:"filePath"`
}

type Instance struct {
	FilePath string `json:"filePath"`
}

func NewCfg() *File {
	return &File{
		Templates:  make([]Template, 0),
		Instances:  make(map[string]([]Instance)),
		Singletons: make([]Singleton, 0),
		NoGit:      false,
	}
}

// everything in .cfg/ (including .cfg/key.json) is not tracked by git
type Key struct {
	Email    string            `json:"email"`
	Remotes  map[string]string `json:"remotes"`
	BaseName string            `json:"baseName"`
}

func NewKey(baseName string) *Key {
	return &Key{
		Email:    viper.GetString("email"),
		Remotes:  make(map[string]string),
		BaseName: baseName,
	}
}

type Snapshot struct {
	Name string `json:"name"`
}

// everything in .cfg/ (including .cfg/snaps) is not tracked by git
// however, snaps are pushed to Confbase servers
type Snaps struct {
	Current   Snapshot   `json:"current"`
	Snapshots []Snapshot `json:"snapshots"`
}

func NewSnaps() *Snaps {
	return &Snaps{
		Current:   Snapshot{Name: "master"},
		Snapshots: []Snapshot{{Name: "master"}},
	}
}

func MustLoadCfg() *File {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	filePath := path.Join(cwd, FileName)

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
	filePath := path.Join(cwd, FileName)

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

func MustLoadKey() *Key {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	keyPath := path.Join(cwd, Dirname, KeyfileName)

	f, err := os.OpenFile(keyPath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to open %v\n", keyPath)
		os.Exit(1)
	}

	key := Key{}
	if err := json.NewDecoder(f).Decode(&key); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse %v\n", keyPath)
		os.Exit(1)
	}
	return &key
}

func (k *Key) MustSerialize(tx *rollback.Tx) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}

	dirPath := path.Join(cwd, Dirname)

	// mkdir if not exist
	_, err = os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to create directory %v\n", dirPath)
			if tx != nil {
				tx.MustRollback()
			}
			os.Exit(1)
		}
		tx.DirsCreated = append(tx.DirsCreated, dirPath)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to stat %v\n", dirPath)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}

	keyBytes, err := json.Marshal(k)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal key\n")
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}

	keyPath := path.Join(dirPath, KeyfileName)

	isCreated := false
	_, err = os.Stat(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			isCreated = true
		} else {
			fmt.Fprintf(os.Stderr, "error: failed to stat %v\n", keyPath)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			if tx != nil {
				tx.MustRollback()
			}
			os.Exit(1)
		}
	}

	if err := ioutil.WriteFile(keyPath, keyBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write %v\n", keyPath)
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, keyPath)
	}
}

func MustLoadSnaps() *Snaps {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	snapsPath := path.Join(cwd, Dirname, SnapsFileName)

	f, err := os.OpenFile(snapsPath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to open %v\n", snapsPath)
		os.Exit(1)
	}

	snaps := Snaps{}
	if err := json.NewDecoder(f).Decode(&snaps); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse %v\n", snapsPath)
		os.Exit(1)
	}
	return &snaps
}

func (s *Snaps) MustSerialize(tx *rollback.Tx) {
	if err := s.Serialize(tx); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func (s *Snaps) Serialize(tx *rollback.Tx) error {
	cwd, err := os.Getwd()
	if err != nil {
		composedErr := fmt.Errorf("error: failed to get working directory\n%v\n", err)
		var txErr error
		if tx != nil {
			txErr = tx.Rollback()
		}
		if txErr != nil {
			composedErr = fmt.Errorf("%v%v\n", composedErr, txErr)
		}
		return composedErr
	}

	dirPath := path.Join(cwd, Dirname)

	// mkdir if not exist
	_, err = os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			composedErr := fmt.Errorf("error: failed to create directory %v\n", dirPath)
			var txErr error
			if tx != nil {
				txErr = tx.Rollback()
			}
			if txErr != nil {
				composedErr = fmt.Errorf("%v%v\n", composedErr, txErr)
			}
			return composedErr
		}
		tx.DirsCreated = append(tx.DirsCreated, dirPath)
	} else if err != nil {
		composedErr := fmt.Errorf("error: failed to stat %v\n", dirPath)
		composedErr = fmt.Errorf("%v%v\n", composedErr, err)
		var txErr error
		if tx != nil {
			txErr = tx.Rollback()
		}
		if txErr != nil {
			composedErr = fmt.Errorf("%v%v\n", composedErr, txErr)
		}
		return composedErr
	}

	snapsBytes, err := json.Marshal(s)
	if err != nil {
		composedErr := fmt.Errorf("error: failed to marshal snapshots file\n")
		var txErr error
		if tx != nil {
			txErr = tx.Rollback()
		}
		if txErr != nil {
			composedErr = fmt.Errorf("%v%v\n", composedErr, txErr)
		}
		return composedErr
	}

	snapsPath := path.Join(dirPath, SnapsFileName)

	isCreated := false
	_, err = os.Stat(snapsPath)
	if err != nil {
		if os.IsNotExist(err) {
			isCreated = true
		} else {
			composedErr := fmt.Errorf("error: failed to stat %v\n", snapsPath)
			composedErr = fmt.Errorf("%v%v\n", composedErr, err)
			var txErr error
			if tx != nil {
				txErr = tx.Rollback()
			}
			if txErr != nil {
				composedErr = fmt.Errorf("%v%v\n", composedErr, txErr)
			}
			return composedErr
		}
	}

	if err := ioutil.WriteFile(snapsPath, snapsBytes, 0644); err != nil {
		composedErr := fmt.Errorf("error: failed to write %v\n", snapsPath)
		var txErr error
		if tx != nil {
			txErr = tx.Rollback()
		}
		if txErr != nil {
			composedErr = fmt.Errorf("%v%v\n", composedErr, txErr)
		}
		return composedErr
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, snapsPath)
	}
	return nil
}

func (cfg *File) Equals(o *File) bool {
	if cfg.NoGit != o.NoGit {
		return false
	}

	if len(cfg.Templates) != len(o.Templates) {
		return false
	}
	for i, v := range cfg.Templates {
		if v != o.Templates[i] {
			return false
		}
	}

	if len(cfg.Instances) != len(o.Instances) {
		return false
	}
	for k, v := range cfg.Instances {
		if len(v) != len(o.Instances[k]) {
			return false
		}
		for i, vv := range v {
			if vv != o.Instances[k][i] {
				return false
			}
		}
	}

	if len(cfg.Singletons) != len(o.Singletons) {
		return false
	}
	for i, v := range cfg.Singletons {
		if v != o.Singletons[i] {
			return false
		}
	}
	return true
}

func mustStage(filePath string) {
	cmd := exec.Command("git", "add", filePath)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to stage %v\n", filePath)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func mustCommit(msg string) {
	cmd := exec.Command("git", "commit", "-m", msg)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to commit\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func (cfg *File) MustStage() {
	for _, t := range cfg.Templates {
		mustStage(t.FilePath)
	}
	for _, t := range cfg.Instances {
		for _, i := range t {
			mustStage(i.FilePath)
		}
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

func (cfg *File) MustPushRaw() {

}

func (s *Snaps) Push() error {
	return nil
}
