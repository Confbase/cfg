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
	cfgFile := dotcfg.MustLoadCfg()
	keyFile := dotcfg.MustLoadKey()
	snapsFile := dotcfg.MustLoadSnaps()

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
		if cfg.Snapshot == "" {
			snapName := snapsFile.Current.Name
			out, err := exec.Command("git", "push", remote, snapName).CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "'git push %v %v' failed\n", remote, snapName)
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				fmt.Fprintf(os.Stderr, "output: %v\n", string(out))
				os.Exit(1)
			}
			fmt.Printf(string(out))
			os.Exit(0)
		}
		out, err := exec.Command("git", "push", remote, cfg.Snapshot).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "'git push %v %v' failed\n", remote, cfg.Snapshot)
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			fmt.Fprintf(os.Stderr, "output: %v\n", string(out))
			os.Exit(1)
		}
		fmt.Printf(string(out))
		os.Exit(0)
	}

	snapName := cfg.Snapshot
	if cfg.Snapshot == "" {
		snapName = snapsFile.Current.Name
	} else if cfg.Snapshot != snapsFile.Current.Name {
		fmt.Fprintf(os.Stderr, "error: pushing a snap other than the current one is not allowed with --no-git\n")
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
	for _, insts := range cfgFile.Instances {
		for _, i := range insts {
			if _, ok := fileSet[i.FilePath]; !ok {
				files = append(files, i.FilePath)
				fileSet[i.FilePath] = true
			}
		}
	}
	for _, s := range cfgFile.Singletons {
		if _, ok := fileSet[s.FilePath]; !ok {
			files = append(files, s.FilePath)
			fileSet[s.FilePath] = true
		}
	}

	fmt.Println("building snapshot...")
	r, w := io.Pipe()
	go func() {
		for _, filePath := range files {
			if err := build.BuildSnap(w, filePath); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		}
		if err := w.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("pushing snapshot...")
	// TODO: some nice ncurses progress thing
	if err := send.SendSnap(remoteValue, r, keyFile.BaseName, snapName, true); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("successfully pushed '%v/%v' to '%v'\n", keyFile.BaseName, snapName, remote)
}
