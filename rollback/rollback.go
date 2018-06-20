package rollback

import (
	"fmt"
	"os"
)

type Tx struct {
	DirsCreated  []string
	FilesCreated []string
}

func NewTx() *Tx {
	return &Tx{
		DirsCreated:  make([]string, 0),
		FilesCreated: make([]string, 0),
	}
}

func (tx *Tx) MustRollback() {
	if err := tx.Rollback(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (tx *Tx) Rollback() error {
	for _, f := range tx.FilesCreated {
		if err := os.Remove(f); err != nil {
			txErr := fmt.Errorf("failed to remove %v\n", f)
			txErr = fmt.Errorf("%v%v\n", txErr, err)
			txErr = fmt.Errorf("%vfailed to rollback changes\n", txErr)
			return txErr
		}
	}
	for _, d := range tx.DirsCreated {
		if err := os.RemoveAll(d); err != nil {
			txErr := fmt.Errorf("failed to remove %v\n", d)
			txErr = fmt.Errorf("%v%v\n", txErr, err)
			txErr = fmt.Errorf("%vfailed to rollback changes\n", txErr)
			return txErr
		}
	}
	return nil
}

func MergeTxErr(err, txErr error) error {
	err = fmt.Errorf("during this error:\n%v", err)
	return fmt.Errorf("%v\ntransaction rollback failed:\n%v", err, txErr)
}
