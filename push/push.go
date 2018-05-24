package push

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Confbase/cfg/dotcfg"
)

func Push(cfg Config) {
	cfgFile := dotcfg.MustLoadCfg()
	keyFile := dotcfg.MustLoadKey()
	snapsFile := dotcfg.MustLoadSnaps()

	remote := ""
	if cfg.Remote == "" {
		remote = "origin"
	} else {
		remote = cfg.Remote
	}
	if _, ok := keyFile.Remotes[remote]; !ok {
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

	snapshots := make([]string, 0)
	if cfg.Snapshot == "" {
		for _, snap := range snapsFile.Snapshots {
			snapshots = append(snapshots, snap.Name)
		}
	} else {
		isValidSnap := false
		for _, snap := range snapsFile.Snapshots {
			if snap.Name == cfg.Snapshot {
				isValidSnap = true
				break
			}
		}
		if !isValidSnap {
			fmt.Fprintf(os.Stderr, "error: %v is not a snapshot\n", cfg.Snapshot)
			os.Exit(1)
		}
		snapshots = append(snapshots, cfg.Snapshot)
	}

	for _, snapName := range snapshots {
		fmt.Printf("pushing snapshot %v to remote %v\n", snapName, remote)
	}
}
