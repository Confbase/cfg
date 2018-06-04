package dotcfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Confbase/cfg/rollback"
)

func NewSnaps() *Snaps {
	return &Snaps{
		Current:   Snapshot{Name: "master"},
		Snapshots: []Snapshot{{Name: "master"}},
	}
}

func MustLoadSnaps() *Snaps {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	snapsPath := filepath.Join(cwd, DirName, SnapsFileName)

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

	dirPath := filepath.Join(cwd, DirName)

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

	snapsPath := filepath.Join(dirPath, SnapsFileName)

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
