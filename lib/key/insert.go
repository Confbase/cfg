package key

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Confbase/cfg/lib/dotcfg"
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

	keyfile := dotcfg.MustLoadKey()
	keyfile.Key = key
	keyfile.MustSerialize(nil)
	fmt.Println("successfully inserted key")
}
