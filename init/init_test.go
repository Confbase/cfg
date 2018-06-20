package init

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdir := filepath.Join(wd, "testdata")
	if err := os.Mkdir(testdir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(testdir); err != nil {
		t.Fatal(err)
	}

	t.Run("--force", func(t *testing.T) {
		cfg := Config{
			AppendGitIgnore:    false,
			OverwriteGitIgnore: false,
			NoGit:              false,
			NoModGitIgnore:     false,
			Force:              true,
		}
		Init(&cfg)
	})

	if err := os.Chdir(wd); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(testdir); err != nil {
		t.Fatal(err)
	}
}
