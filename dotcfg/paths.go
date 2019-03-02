package dotcfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetRelativeToBaseDir takes a baseDir and a filePath.
// The relative path from baseDir to filePath is returned.
// If filePath stats with /, it is interpreted as an absolute path.
// Otherwise, it is is interpreted as a relative path.
func GetRelativeToBaseDir(baseDir, filePath string) (string, error) {
	if len(baseDir) == 0 {
		return "", fmt.Errorf("base directory is an empty string")
	}
	if len(filePath) == 0 {
		return "", fmt.Errorf("file path is an empty string")
	}
	if filePath[0] != '/' {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		filePath = filepath.Join(cwd, filePath)
	} else {
		filePath = filepath.Clean(filePath)
	}
	baseDir = filepath.Clean(baseDir)
	if !strings.HasPrefix(filePath, baseDir) {
		i := 0
		for {
			if baseDir[i] != filePath[i] {
				break
			}
			i++
		}
		nonDotDotPart := filePath[i:len(filePath)]
		numSlashes := strings.Count(baseDir[i:len(baseDir)], "/")
		relPath := ""
		for j := 0; j < numSlashes+1; j++ {
			relPath = relPath + "../"
		}
		return relPath + nonDotDotPart, nil
	}
	if filePath == baseDir {
		return ".", nil
	}
	return filePath[len(baseDir)+1 : len(filePath)], nil
}

// GetAbsAndRelPaths takes a baseDir and an absolute or relative (to cwd) file
// path. The absolute file path and relative (to baseDir) file paths are
// returned, along with any errors which occured.
func GetAbsAndRelPaths(baseDir, filePath string) (string, string, error) {
	relPath, err := GetRelativeToBaseDir(baseDir, filePath)
	return filepath.Join(baseDir, relPath), relPath, err
}
