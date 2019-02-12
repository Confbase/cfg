package rm

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"

	initPkg "github.com/Confbase/cfg/init"
	"github.com/Confbase/cfg/mark"
)

func TestRm(t *testing.T) {
	// creating a temporary directory in /tmp in case os.exit gets called and
	// we don't have a chance to cleanup
	testdir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatal(err)
	}

	nestedTestdir, err := ioutil.TempDir(testdir, "nested")
	if err != nil {
		t.Fatal(err)
	}

	// automatic cleanup of testing directory
	// NOTE: does not work if os.exit is called.
	defer os.RemoveAll(nestedTestdir)
	defer os.RemoveAll(testdir)

	if err := os.Chdir(testdir); err != nil {
		t.Fatal(err)
	}

	testFilePaths := make([]string, 2)
	filePath1, err := ioutil.TempFile(testdir, "*.json")
	if err != nil {
		t.Fatal(err)
	} else {
		testFilePaths[0] = filePath1.Name()
		w := bufio.NewWriter(filePath1)
		_, err := w.WriteString("{\"test_key\": 6}\n")
		if err != nil {
			t.Fatal(err)
		}

		w.Flush()
	}

	filePath2, err := ioutil.TempFile(nestedTestdir, "*.json")
	if err != nil {
		t.Fatal(err)
	} else {
		testFilePaths[1] = filePath2.Name()
		w := bufio.NewWriter(filePath2)

		_, err := w.WriteString("{\"test_key\": 5}\n")
		if err != nil {
			t.Fatal(err)
		}

		w.Flush()
	}

	initCfg := initPkg.Config{
		AppendGitIgnore:    false,
		OverwriteGitIgnore: false,
		NoGit:              false,
		NoModGitIgnore:     false,
		Force:              true,
		Dest:               testdir,
	}

	markCfg1 := mark.Config{
		Template: "test_instance1",
		Targets:  testFilePaths[:1],
	}

	markCfg2 := mark.Config{
		Template: "test_instance2",
		Targets:  testFilePaths[1:],
	}

	rmCfg1 := Config{
		ToRemove:  markCfg1.Targets,
		Recursive: false,
	}

	rmCfg2 := Config{
		ToRemove:  markCfg2.Targets,
		Recursive: false,
	}

	t.Run("BasicRmTest", func(t *testing.T) {
		if err := initPkg.Init(&initCfg); err != nil {
			t.Fatal(err)
		}

		mark.Mark(&markCfg1)
		mark.Mark(&markCfg2)

		if err := Rm(&rmCfg1); err != nil {
			t.Fatal(err)
		}

		if err := Rm(&rmCfg2); err != nil {
			t.Fatal(err)
		}
	})
}
