package dotcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Confbase/cfg/lib/rollback"
	"github.com/Confbase/cfg/lib/util"
)

const (
	Dirname     = ".cfg"      // this dir resides in ./
	FileName    = ".cfg.json" // this file resides in ./
	KeyfileName = "key.json"  // this file resides in ./.cfg/
)

type Template struct {
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
	FileType string `json:"fileType"`
}

// .cfg.json is tracked by git
type File struct {
	Templates  []Template            `json:"templates"`
	Instances  map[string]([]string) `json:"instances"`
	Singletons []string              `json:"singletons"`
}

// .cfg/ (including .cfg/key.json) is not tracked by git
type Key struct {
	Email string `json:"email"`
	Key   string `json:"key"`

	EntryPoint string            `json:"entryPoint"` // Confbase API base URL
	Remotes    map[string]string `json:"remotes"`
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
	isExist, err := util.Exists(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to stat %v\n", filePath)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}
	if !isExist {
		isCreated = true
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
	isExist, err := util.Exists(keyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to stat %v\n", keyPath)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if tx != nil {
			tx.MustRollback()
		}
		os.Exit(1)
	}
	if !isExist {
		isCreated = true
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

func (t *Template) Equals(o *Template) bool {
	if t.Name != o.Name {
		return false
	}
	if t.FilePath != o.FilePath {
		return false
	}
	if t.FileType != o.FileType {
		return false
	}
	return true
}

func (cfg *File) Equals(o *File) bool {
	if len(cfg.Templates) != len(o.Templates) {
		return false
	}
	for i, v := range cfg.Templates {
		if !v.Equals(&o.Templates[i]) {
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
