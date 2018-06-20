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
	if err := k.Serialize(tx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (k *Key) Serialize(tx *rollback.Tx) error {
	cwd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("failed to get working directory\n%v", err)
		if tx != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
		}
		return err
	}

	dirPath := filepath.Join(cwd, DirName)

	// mkdir if not exist
	_, err = os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			err = fmt.Errorf("failed to create directory %v\n%v", dirPath, err)
			if tx != nil {
				if txErr := tx.Rollback(); txErr != nil {
					err = rollback.MergeTxErr(err, txErr)
				}
			}
			return err
		}
		tx.DirsCreated = append(tx.DirsCreated, dirPath)
	} else if err != nil {
		err = fmt.Errorf("failed to stat %v\n%v", dirPath, err)
		if tx != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
		}
		return err
	}

	keyBytes, err := json.Marshal(k)
	if err != nil {
		err = fmt.Errorf("failed to marshal key\n%v", err)
		if tx != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
		}
		return err
	}

	keyPath := filepath.Join(dirPath, KeyfileName)

	isCreated := false
	_, err = os.Stat(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			isCreated = true
		} else {
			err = fmt.Errorf("failed to stat %v\n%v", keyPath, err)
			if tx != nil {
				if txErr := tx.Rollback(); txErr != nil {
					err = rollback.MergeTxErr(err, txErr)
				}
			}
			return err
		}
	}

	if err := ioutil.WriteFile(keyPath, keyBytes, 0644); err != nil {
		err = fmt.Errorf("failed to write %v\n%v", keyPath, err)
		if tx != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
		}
		return err
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, keyPath)
	}
	return nil
}
