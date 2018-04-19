package snap

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Confbase/cfg/lib/dotcfg"
)

func mustUnCheckout() {
	coPrevCmd := exec.Command("git", "checkout", "-")
	if err := coPrevCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to checkout previous branch\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func New(name string) {

	snaps := dotcfg.MustLoadSnaps()
	for _, s := range snaps.Snapshots {
		if s.Name == name {
			fmt.Fprintf(os.Stderr, "error: a snapshot named '%v' already exists\n", name)
			os.Exit(1)
		}
	}

	cfg := dotcfg.MustLoadCfg()
	if !cfg.NoGit {
		stsCmd := exec.Command("git", "status", "-s")
		stsBytes, stsErr := stsCmd.Output()
		if stsErr != nil {
			fmt.Fprintf(os.Stderr, "error: failed to run 'git status -s'\n%v", stsErr)
			os.Exit(1)
		}
		sts := string(stsBytes)
		if sts != "" {
			for _, line := range strings.Split(sts, "\n") {
				if len(line) >= 2 && line[0] == '?' && line[1] == '?' {
					fmt.Fprintf(os.Stderr, "error: found untracked files; you must\n\n")
					fmt.Fprintf(os.Stderr, "    1. track these files via 'cfg track', 'cfg tag', or 'cfg mark'\n")
					fmt.Fprintf(os.Stderr, "    or 2. add these files to the .gitignore\n")
					fmt.Fprintf(os.Stderr, "    or 3. stash these files with 'git stash'\n")
					fmt.Fprintf(os.Stderr, "    or 4. delete these files\n\n")
					fmt.Fprintf(os.Stderr, "'git status -s' output (?? indicates a file is untracked):\n%v", sts)
					os.Exit(1)
				}
			}
			fmt.Printf("committing changes before creating new snapshot...")
			cfg.MustStage()
			cfg.MustCommit()
			fmt.Printf("OK\n")
		}

		coNewCmd := exec.Command("git", "checkout", "-b", name)
		if err := coNewCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to checkout new branch '%v' for snapshot\n", name)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("switched to new branch '%v'", name)

		// TODO: push to Confbase git servers
	}

	snaps.Snapshots = append(snaps.Snapshots, dotcfg.Snapshot{Name: name})
	snaps.Current = name
	if err := snaps.Serialize(nil); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to serialize snapshots file\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)

		if !cfg.NoGit {
			mustUnCheckout()
		}
		os.Exit(1)
	}
	if err := snaps.Push(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to push snapshots file\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)

		if !cfg.NoGit {
			mustUnCheckout()
		}
		os.Exit(1)
	}

	cfg.MustPushRaw()
}
