package dotcfg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestMustLoadCfg(t *testing.T) {
	const (
		crashEnvKey = "MUST_LOAD_CFG_SHOULD_CRASH"
		crashEnvVal = "true"

		destFileName = ".cfg.json"
		srcFileName  = "TestMustLoadCfg.cfg.json"
	)

	if os.Getenv(crashEnvKey) == crashEnvVal {
		MustLoadCfg()
		return
	}

	_, err := os.Stat(destFileName)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("%v already exists", destFileName)
	}

	destFile, err := os.Create(destFileName)
	if err != nil {
		t.Fatalf("failed to create %v\nerror: %v", destFileName, err)
	}
	defer destFile.Close()

	srcFile, err := os.Open(srcFileName)
	if err != nil {
		t.Fatalf("failed to open %v\nerror: %v", srcFileName, err)
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		t.Fatalf("failed to copy %v to %v\nerror: %v", srcFileName, destFileName, err)
	}

	expect := NewCfg()
	expect.NoGit = true
	expect.Templates = append(
		expect.Templates,
		Template{
			Name:     "test MustLoadCfg templ name 0",
			FilePath: ".TestMustLoadCfg_0.json",
		},
		Template{
			Name:     "test MustLoadCfg templ name 1",
			FilePath: ".TestMustLoadCfg_1.json",
		},
	)
	expect.Instances["test MustLoadCfg templ name 0"] = make([]Instance, 0)
	expect.Instances["test MustLoadCfg templ name 0"] = append(
		expect.Instances["test MustLoadCfg instance key"],
		Instance{FilePath: ".test_MustLoadCfg_instance_element_0"},
		Instance{FilePath: ".test_MustLoadCfg_instance_element_1"},
	)
	expect.Instances["test MustLoadCfg templ name 1"] = make([]Instance, 0)
	expect.Instances["test MustLoadCfg templ name 1"] = append(
		expect.Instances["test MustLoadCfg instance key"],
		Instance{FilePath: ".test_MustLoadCfg_instance_element_2"},
		Instance{FilePath: ".test_MustLoadCfg_instance_element_3"},
	)
	expect.Singletons = make([]Singleton, 0)
	expect.Singletons = append(
		expect.Singletons,
		Singleton{FilePath: ".test_MustLoadCfg_singletons_element_0"},
		Singleton{FilePath: ".test_MustLoadCfg_singletons_element_1"},
	)

	got := MustLoadCfg()
	if !got.Equals(expect) {
		t.Errorf("expected %v but got %v", expect, got)
	}

	if err := os.Remove(destFileName); err != nil {
		t.Fatalf("failed to remove %v", destFileName)
	}

	crashCmd := exec.Command(os.Args[0], "-test.run=TestMustLoadCfg")
	crashCmd.Env = append(os.Environ(), fmt.Sprintf("%v=%v", crashEnvKey, crashEnvVal))
	err = crashCmd.Run()
	if err == nil {
		t.Fatalf("expected non-zero exit status; got zero")
	}
	e, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("failed to cast crashCmd err to *exec.ExitError")
	}
	if e.Success() {
		t.Fatalf("expected non-zero exit status; got zero")
	}
}
