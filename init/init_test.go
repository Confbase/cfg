package init

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdir := filepath.Join(wd, "testdir")
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

		if err := Init(&cfg); err != nil {
			t.Fatal(err)
		}
	})

	if err := os.Chdir(wd); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(testdir); err != nil {
		t.Fatal(err)
	}
}

// Tests that you can call cfg init /some/other/path/far/far/away
func TestInitOverThere(t *testing.T) {
	// creating a temporary directory in /tmp in case os.exit gets called and
	// we don't have a chance to cleanup
	testdir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatal(err)
	}

	// automatic cleanup of testing directory
	// NOTE: does not work if os.exit is called.
	defer os.RemoveAll(testdir)

	t.Run("--force", func(t *testing.T) {
		defer os.RemoveAll(testdir)
		cfg := Config{
			AppendGitIgnore:    false,
			OverwriteGitIgnore: false,
			NoGit:              false,
			NoModGitIgnore:     false,
			Force:              true,
			Dest:               testdir,
		}

		if err := Init(&cfg); err != nil {
			t.Fatal(err)
		}
	})
}
