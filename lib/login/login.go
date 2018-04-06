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

func Login() {
	fmt.Printf("email: ")
	var email string
	if _, err := fmt.Scanln(&email); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("password: ")
	passBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	password := string(passBytes[:])
	fmt.Printf("\n")

	uri := fmt.Sprintf("%s/api/auth/verify", viper.GetString("entryPoint"))

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
