package rm

import (
	"fmt"
	"os"

	"github.com/Confbase/cfg/unmark"
)

// func MustRm(filePaths []string) {
// 	if err := Rm(filePaths); err != nil {
// 		fmt.Fprintf(os.Stderr, "error: %v\n", err)
// 		os.Exit(1)
// 	}
// }

// Main entry point for cfg rm.
func MustRm(cfg *Config) {
	if err := Rm(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// CollectFiles takes in a list of filepaths, and
// recursively walks down the tree for each filepath,
// eventually returning all files contained
// to be removed
// func CollectFiles(filePaths []string) []string, error {
// 	// Poor man's proxy for a hashset
// 	set := make(map[string]struct{})

// 	for _, filePath := range filePaths {
// 		err := filepath.Walk(filepath, func(path string, info os.FileInfo, err error) error {

// 		})
// 	}
// }

// RmRecursive recursively calls rm on all children
// of the directory given by filePath
// func RmRecursive(filePath string) error {

// }

// Rm removes specified files
// First, we unmark the files, and then we remove them
// from the local filesystem
// func Rm(cfg *Config) {

// }

func Rm(cfg *Config) error {
	// TODO: intelligently order filePaths to remove
	// all instances of a templ before removing the templ
	if err := unmark.Unmark(cfg.ToRemove); err != nil {
		return err
	}
	for _, filePath := range cfg.ToRemove {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return nil
}
