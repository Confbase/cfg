package dotcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	Dirname     = ".cfg"      // this dir resides in ./
	FileName    = ".cfg.json" // this file resides in ./
	KeyfileName = "key.json"  // this file resides in ./.cfg/
)

// .cfg.json is tracked by git
type File struct {
	Templates []string          `json:"templates"`
	Instances map[string]string `json:"instances"`
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

func (f *File) MustSerialize() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	filePath := path.Join(cwd, FileName)

	cfgBytes, err := json.Marshal(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal key\n")
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filePath, cfgBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write %v\n", filePath)
		os.Exit(1)
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

func (k *Key) MustSerialize() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}

	dirPath := path.Join(cwd, Dirname)

	// mkdir if not exist
	_, err = os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to create directory %v\n", dirPath)
			os.Exit(1)
		}
	}

	keyBytes, err := json.Marshal(k)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal key\n")
		os.Exit(1)
	}

	keyPath := path.Join(dirPath, KeyfileName)
	if err := ioutil.WriteFile(keyPath, keyBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write %v\n", keyPath)
		os.Exit(1)
	}
}
