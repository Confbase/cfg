package remote

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Confbase/cfg/dotcfg"
)

func Add(name, url string) {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			fmt.Fprintf(os.Stderr, "error: remote name must be alphanumeric\n")
			os.Exit(1)
		}
	}

	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	key := dotcfg.MustLoadKey(baseDir)
	if _, ok := key.Remotes[name]; ok {
		fmt.Fprintf(os.Stderr, "error: a remote named %v already exists\n", name)
		os.Exit(1)
	}
	key.Remotes[name] = url

	cfgFile := dotcfg.MustLoadCfg(baseDir)
	if !cfgFile.NoGit {
		out, err := exec.Command("git", "remote", "add", name, url).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "'git remote add %v %v' failed\n", name, url)
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			fmt.Fprintf(os.Stderr, "output: %v\n", string(out))
			os.Exit(1)
		}
	}

	key.MustSerialize(baseDir, nil)
}

func Remove(name string) {
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	key := dotcfg.MustLoadKey(baseDir)
	if _, ok := key.Remotes[name]; !ok {
		fmt.Fprintf(os.Stderr, "error: there is no remote named %v\n", name)
		os.Exit(1)
	}
	delete(key.Remotes, name)

	cfgFile := dotcfg.MustLoadCfg(baseDir)
	if !cfgFile.NoGit {
		out, err := exec.Command("git", "remote", "rm", name).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "'git remote rm %v' failed\n", name)
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			fmt.Fprintf(os.Stderr, "output: %v\n", string(out))

			// serialize the key file anyways
			// so the remote is removed
			key.MustSerialize(baseDir, nil)

			os.Exit(1)
		}
	}

	key.MustSerialize(baseDir, nil)
}

func Ls() {
	baseDir, err := dotcfg.GetBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	key := dotcfg.MustLoadKey(baseDir)
	for name, url := range key.Remotes {
		fmt.Printf("%v %v\n", name, url)
	}
}

func ExtractTeam(remote string) (string, bool) {
	if !strings.HasPrefix(remote, "https://") {
		return "", false
	}
	dotIndex := strings.Index(remote, ".")
	if dotIndex == -1 {
		return "", false
	}
	return remote[len("https://"):dotIndex], true
}

func ExtractBase(remote string) (string, bool) {
	if !strings.HasPrefix(remote, "https://") {
		return "", false
	}

	dotIndex := strings.Index(remote, ".")
	if dotIndex == -1 || dotIndex == len(remote)-1 {
		return "", false
	}

	relSlashIndex := strings.Index(remote[dotIndex+1:len(remote)], "/")
	if relSlashIndex == -1 {
		return "", false
	}
	slashIndex := relSlashIndex + dotIndex + 1
	if slashIndex == len(remote)-1 {
		return "", false
	}

	return remote[slashIndex+1 : len(remote)], true
}
