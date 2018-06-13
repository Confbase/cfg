package rm

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/unmark"
)

func MustRm(filePaths []string) {
	if err := Rm(filePaths); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Rm(filePaths []string) error {
	// TODO: intelligently order filePaths to remove
	// all instances of a templ before removing the templ
	if err := unmark.Unmark(filePaths); err != nil {
		return err
	}
	for _, filePath := range filePaths {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return nil
}
