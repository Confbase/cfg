package fetch

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/viper"
)

// myteam:mybase/config.yml

func argsFromTarget(target string) (string, string, string) {
	var team, base, filePath string

	firstColon := strings.Index(target, ":")
	if firstColon == -1 || firstColon == len(target)-1 {
		fmt.Fprintf(os.Stderr, "error: failed to extract team from target\n")
		fmt.Fprintf(os.Stderr, "target must be of the form\nmyteam:mybase/path/to/myfile.yml")
		os.Exit(1)
	}
	team = target[0:firstColon]

	firstSlash := strings.Index(target, "/")
	if firstSlash == -1 || firstSlash == len(target)-1 {
		fmt.Fprintf(os.Stderr, "error: failed to extract team from target\n")
		fmt.Fprintf(os.Stderr, "target must be of the form\nmyteam:mybase/path/to/myfile.yml")
		os.Exit(1)
	}
	base = target[firstColon+1 : firstSlash]
	filePath = target[firstSlash+1 : len(target)]

	return team, base, filePath
}

func Fetch(snapshot, target, email, password string, noGlobal bool) {
	team, base, filePath := argsFromTarget(target)

	// TODO: cleanly get entryPoint
	// - global config won't exist initially by default
	// - could use remotes from key file
	// - following code doesn't account for `noGlobal` var
	entryPoint := viper.GetString("entryPoint")
	if entryPoint == "" {
		fmt.Fprintf(os.Stderr, "error: failed to get entry point from global config file\n")
		os.Exit(1)
	}

	if email == "" {
		if !noGlobal {
			email = viper.GetString("email")
		}
	}
	if email == "" {
		fmt.Printf("email: ")
		if _, err := fmt.Scanln(&email); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	var passBytes []byte
	if password == "" {
		if !noGlobal {
			password = viper.GetString("password")
		}
	}
	if password == "" {
		fmt.Printf("password: ")
		var err error
		passBytes, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		password = string(passBytes[:])
		fmt.Printf("\n")
	} else {
		passBytes = []byte(password)
	}

	uri := fmt.Sprintf("http://%v.%v/%v/raw/%v/%v", team, entryPoint, base, snapshot, filePath)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create request for '%v'\n%v", filePath, err)
		os.Exit(1)
	}
	hashedBytes := md5.Sum(passBytes)
	hashedPass := fmt.Sprintf("%x", hashedBytes)
	req.SetBasicAuth(email, hashedPass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to fetch file '%v'\n%v", filePath, err)
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Fprintf(os.Stderr, "error: invalid credentials while attempting to fetch '%v'\n", target)
		os.Exit(1)
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Fprintf(os.Stderr, "error: '%v' does not exist\n", target)
		fmt.Fprintf(os.Stderr, "uri: '%v'\n", uri)
		os.Exit(1)
	} else if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error: received non-200 status code when fetching '%v'\n", target)
		os.Exit(1)
	}

	// TODO: open file name instead of filePath
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to open file '%v'\n%v", filePath, err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to fetch file '%v'\n%v", filePath, err)
		os.Exit(1)
	}
}
