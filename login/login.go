package login

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/viper"
)

func Login(email, password string) {
	if email == "" {
		fmt.Printf("email: ")
		if _, err := fmt.Scanln(&email); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	var passBytes []byte
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

	uri := fmt.Sprintf("http://%s/api/auth/verify", viper.GetString("entryPoint"))

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create GET request\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	hashedBytes := md5.Sum(passBytes)
	hashedPass := fmt.Sprintf("%x", hashedBytes)
	req.SetBasicAuth(email, hashedPass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to send GET request\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Fprintf(os.Stderr, "invalid credentials; keeping original config file")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "error: received unexpected status code %v\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Printf("success! writing to %v\n", viper.ConfigFileUsed())

	viper.Set("email", email)
	viper.Set("password", password)

	viper.WriteConfig()
}
