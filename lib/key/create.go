package key

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/confbase/cfg/lib/dotcfg"
	"github.com/confbase/cfg/lib/remote"
)

func Create(team, base, name string, canRead, canWrite bool, expiry int) {
	email := viper.GetString("email")
	password := viper.GetString("password")
	if email == "" || password == "" {
		fmt.Fprintf(os.Stderr, "credentials not found in %v; ", viper.ConfigFileUsed())
		fmt.Fprintf(os.Stderr, "you must login to use \"cfg key create\"\n")
		os.Exit(1)
	}

	key := dotcfg.MustLoadKey()

	if team == "" {
		r, ok := key.Remotes["origin"]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: you must specify --team or add an \"origin\" remote\n")
			os.Exit(1)
		}
		team, ok = remote.ExtractTeam(r)
		if !ok {
			fmt.Fprintf(os.Stderr, "error: failed to parse team from \"origin\" remote\n")
			fmt.Fprintf(os.Stderr, "you must either\n\n")
			fmt.Fprintf(os.Stderr, "1. specify --team, or\n")
			fmt.Fprintf(os.Stderr, "2. add a remote called \"origin\" of the form (https://myteam.confbase.com/mybase)\n")
			os.Exit(1)
		}
	}
	if base == "" {
		r, ok := key.Remotes["origin"]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: you must specify --base or add an \"origin\" remote\n")
			os.Exit(1)
		}
		base, ok = remote.ExtractBase(r)
		if !ok {
			fmt.Fprintf(os.Stderr, "error: failed to parse base from \"origin\" remote\n")
			fmt.Fprintf(os.Stderr, "you must either\n\n")
			fmt.Fprintf(os.Stderr, "1. specify --base, or\n")
			fmt.Fprintf(os.Stderr, "2. add a remote called \"origin\" of the form (https://myteam.confbase.com/mybase)\n")
			os.Exit(1)
		}
	}

	fmt.Printf("creating key for team \"%v\", base \"%v\", key name \"%v\"\n", team, base, name)

	data := struct {
		Team     string `json:"team"`
		Base     string `json:"base"`
		Name     string `json:"name"`
		CanRead  bool   `json:"canRead"`
		CanWrite bool   `json:"canWrite"`
		Expiry   int    `json:"expiry"`
	}{
		Team:     team,
		Base:     base,
		Name:     name,
		CanRead:  canRead,
		CanWrite: canWrite,
		Expiry:   expiry,
	}

	payload, err := json.Marshal(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal data to JSON\n")
		os.Exit(1)
	}

	uri := fmt.Sprintf("%s/api/auth/keys/generate", viper.GetString("entryPoint"))

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create POST request\n")
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	hashedBytes := md5.Sum([]byte(password))
	hashedPass := fmt.Sprintf("%x", hashedBytes)
	req.SetBasicAuth(email, hashedPass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to send POST request\n")
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Fprintf(os.Stderr, "error: received non-201 status code\n")
		fmt.Fprintf(os.Stderr, "are you sure the specified base exists?\n")
		fmt.Fprintf(os.Stderr, "are you sure the credentials in the global config file are correct?\n")
		os.Exit(1)
	}

	respKey := struct {
		Key string `json:"key"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&respKey); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to deserialize response JSON\n")
		os.Exit(1)
	}
	if respKey.Key == "" {
		fmt.Fprintf(os.Stderr, "error: received empty key\n")
		os.Exit(1)
	}

	key.Key = respKey.Key
	key.MustSerialize()

	fmt.Println("successfully created key!")
}
