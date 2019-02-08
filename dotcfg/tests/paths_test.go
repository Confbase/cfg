package dotcfg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Confbase/cfg/dotcfg"
	"github.com/magiconair/properties/assert"
)

func TestConvertPathToRelative(t *testing.T) {
	tests := []struct {
		cwd          string
		filePath     string
		relativePath string
		err          error
	}{
		{cwd: "/", filePath: "/a", relativePath: "a", err: nil},
		{cwd: "/", filePath: "/a/b", relativePath: "a/b", err: nil},
		{cwd: "/a/b", filePath: "/c", relativePath: "../../c", err: nil},
		{cwd: "/", filePath: "/", relativePath: ".", err: nil},
	}

	for _, test := range tests {
		tempDir := os.TempDir()
		defer os.RemoveAll(tempDir)

		testDir := filepath.Join(tempDir, test.cwd)

		testPath := test.filePath
		if filepath.IsAbs(test.filePath) {
			testPath = filepath.Join(tempDir, test.filePath)
		}

		err := os.MkdirAll(testDir, 0700)
		if err != nil {
			t.Fatal(err)
		}

		baseTestDirs := strings.Split(testDir, "/")
		if len(baseTestDirs) >= 2 {
			defer os.RemoveAll(filepath.Join(tempDir, baseTestDirs[1]))
		}

		err = os.Chdir(testDir)
		if err != nil {
			t.Fatal(err)
		}

		relativePath, err := dotcfg.ConvertPathToRelative(testPath)
		assert.Equal(t, test.relativePath, relativePath)
		assert.Equal(t, test.err, err)
	}
}
