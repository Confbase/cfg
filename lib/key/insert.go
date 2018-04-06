package key

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/confbase/cfg/lib/dotcfg"
)

func Insert() {
	fmt.Printf("paste the key: ")
	keyBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	key := string(keyBytes[:])
	fmt.Printf("\n")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to get working directory\n")
		os.Exit(1)
	}
	keyfilePath := path.Join(cwd, dotcfg.Dirname, dotcfg.KeyfileName)

	f, err := os.OpenFile(keyfilePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to open %v\n", keyfilePath)
		os.Exit(1)
	}
	defer f.Close()

	keyfile := dotcfg.Key{}
	if err := json.NewDecoder(f).Decode(&keyfile); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse %v\n", keyfilePath)
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	f.Close()

	keyfile.Key = key

	kfBytes, err := json.Marshal(keyfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal keyfile\n")
		os.Exit(1)
	}

	if err := ioutil.WriteFile(keyfilePath, kfBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write %v\n", keyfilePath)
		os.Exit(1)
	}

	fmt.Printf("successfully inserted key into %v\n", keyfilePath)
}
