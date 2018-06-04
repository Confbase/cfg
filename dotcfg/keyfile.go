package dotcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/Confbase/cfg/rollback"
)

func NewKey(baseName string) *Key {
	return &Key{
		Email:    viper.GetString("email"),
		Remotes:  make(map[string]string),
		BaseName: baseName,
	}
}

func MustLoadKey() *Key {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	keyPath := filepath.Join(cwd, DirName, KeyfileName)

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

	dirPath := filepath.Join(cwd, DirName)

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

	keyPath := filepath.Join(dirPath, KeyfileName)

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
