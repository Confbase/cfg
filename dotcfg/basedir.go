package dotcfg

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetBaseDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	target := filepath.Join(cwd, FileName)
	for {
		if _, err := os.Stat(target); err != nil && os.IsNotExist(err) {
			// go up one directory
			parentDir := filepath.Dir(filepath.Dir(target))
			if parentDir == "/" {
				return "", fmt.Errorf("no %v found in file path", FileName)
			}
			target = filepath.Join(parentDir, FileName)
		} else if err != nil {
			return "", err
		} else {
			break
		}
	}
	return filepath.Dir(target), nil
}
