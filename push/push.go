package push

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Confbase/cfg/dotcfg"
	"github.com/Confbase/cfgd/cfgsnap/build"
	"github.com/Confbase/cfgd/cfgsnap/send"
)

func Push(cfg Config) {
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	cfgFile := dotcfg.MustLoadCfg(baseDir)
	keyFile := dotcfg.MustLoadKey(baseDir)
	snapsFile := dotcfg.MustLoadSnaps(baseDir)

	if len(keyFile.Remotes) == 0 {
		fmt.Fprintf(os.Stderr, "error: there are no remotes\n")
		fmt.Fprintf(os.Stderr, "a remote can be added with ")
		fmt.Fprintf(os.Stderr, "'cfg remote add <remote-name> <remote-url>'\n")
		os.Exit(1)
	}

	remote := ""
	if cfg.Remote == "" {
		remote = "origin"
	} else {
		remote = cfg.Remote
	}
	remoteValue, ok := keyFile.Remotes[remote]
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %v is not a remote\n", remote)
		os.Exit(1)
	}

	if !cfgFile.NoGit {

		snapName := cfg.Snapshot
		if snapName == "" {
			snapName = snapsFile.Current.Name
		}

		pushCmd := exec.Command("git", "push", remote, snapName)
		if cfg.IsForce {
			pushCmd = exec.Command("git", "push", "--force",
				remote, snapName)
		}
		out, err := pushCmd.CombinedOutput()
		if err != nil {
			if cfg.IsForce {
				fmt.Fprintf(os.Stderr,
					"'git push --force %v %v' failed\n",
					remote, snapName)
			} else {
				fmt.Fprintf(os.Stderr, "'git push %v %v' failed\n",
					remote, snapName)
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			fmt.Fprintf(os.Stderr, "output: %v\n", string(out))
			os.Exit(1)
		}
		fmt.Print(string(out))
		os.Exit(0)
	}

	snapName := cfg.Snapshot
	if cfg.Snapshot == "" {
		snapName = snapsFile.Current.Name
	} else if cfg.Snapshot != snapsFile.Current.Name {
		fmt.Fprintf(os.Stderr, "error: pushing a snap other than "+
			"the current one is not allowed with --no-git\n")
		os.Exit(1)
	}

	fileSet := make(map[string]bool) // set of filepaths
	files := make([]string, 0)
	for _, t := range cfgFile.Templates {
		if _, ok := fileSet[t.FilePath]; !ok {
			files = append(files, t.FilePath)
			fileSet[t.FilePath] = true
		}
	}
	for _, inst := range cfgFile.Instances {
		if _, ok := fileSet[inst.FilePath]; !ok {
			files = append(files, inst.FilePath)
			fileSet[inst.FilePath] = true
		}
	}
	for _, s := range cfgFile.Singletons {
		if _, ok := fileSet[s.FilePath]; !ok {
			files = append(files, s.FilePath)
			fileSet[s.FilePath] = true
		}
	}

	fmt.Println("building snapshot...")
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "warning: no files to push\n")
		os.Exit(0)
	}
	snapRdr, err := build.BuildSnap(files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	for _, filePath := range files {
		r, err := build.BuildSnap(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		snapRdr = io.MultiReader(snapRdr, r)
	}

	fmt.Println("pushing snapshot...")
	// TODO: some nice ncurses progress thing
	if err := send.SendSnap(
		remoteValue,
		snapRdr,
		keyFile.BaseName,
		snapName,
		true,
	); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("successfully pushed '%v/%v' to '%v'\n",
		keyFile.BaseName, snapName, remote)
}
