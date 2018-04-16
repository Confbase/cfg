package remote

import (
	"fmt"
	"os"
	"strings"

	"github.com/Confbase/cfg/lib/dotcfg"
)

func Add(name, url string) {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			fmt.Fprintf(os.Stderr, "error: remote name must be alphanumeric\n")
			os.Exit(1)
		}
	}

	key := dotcfg.MustLoadKey()
	if _, ok := key.Remotes[name]; ok {
		fmt.Fprintf(os.Stderr, "error: a remote named %v already exists\n", name)
		os.Exit(1)
	}
	key.Remotes[name] = url
	key.MustSerialize()
}

func Remove(name string) {
	key := dotcfg.MustLoadKey()
	if _, ok := key.Remotes[name]; !ok {
		fmt.Fprintf(os.Stderr, "error: there is no remote named %v\n", name)
		os.Exit(1)
	}
	delete(key.Remotes, name)
	key.MustSerialize()
}

func Ls() {
	key := dotcfg.MustLoadKey()
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
