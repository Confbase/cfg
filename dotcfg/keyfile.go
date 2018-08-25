package dotcfg

import (
	"encoding/json"
	"fmt"
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

func (k *Key) MustSerialize(baseDir string, tx *rollback.Tx) {
	if err := k.Serialize(baseDir, tx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (k *Key) Serialize(baseDir string, tx *rollback.Tx) error {
	dirPath := filepath.Join(baseDir, DirName)

	// mkdir if not exist
	if _, err := os.Stat(dirPath); err != nil && os.IsNotExist(err) {
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

	keyPath := filepath.Join(dirPath, KeyfileName)

	isCreated := false
	if _, err := os.Stat(keyPath); err != nil {
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

	f, err := os.Create(keyPath)
	if err != nil {
		err = fmt.Errorf("failed to create or open %v\n%v", keyPath, err)
		if tx != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
		}
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(k); err != nil {
		err = fmt.Errorf("failed to write %v\n%v", keyPath, err)
		if tx != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = rollback.MergeTxErr(err, txErr)
			}
		}
		return err
	}

	// explicitly close, to verify the file has been written
	if err := f.Close(); err != nil {
		err = fmt.Errorf("failed to close file %v\n%v", keyPath, err)
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("during error:\n%v\ntransaction rollback failed with error:\n%v", err, txErr)
		}
		return err
	}

	if isCreated {
		tx.FilesCreated = append(tx.FilesCreated, keyPath)
	}
	return nil
}
