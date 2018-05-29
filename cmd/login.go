// Copyright © 2018 Thomas Fischer
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Confbase/cfg/login"
)

var loginUsername, loginPassword string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Store credentials in global config file",
	Long: `WARNING: it is strongly recommended to use per-base access keys instead.

"cfg login" prompts for a username and password and stores these credentials
in plain text in the global config file upon successfully logging in.

See "cfg key" for per-base access key management.`,
	Run: func(cmd *cobra.Command, args []string) {
		login.Login(loginUsername, loginPassword)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "username")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "password")
}
