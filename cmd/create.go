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

	"github.com/confbase/cfg/lib/key"
)

var team, base, name string
var canRead, canWrite bool
var expiry int

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create, insert a key using credentials in ~/.cfg.yaml",
	Long: `This command creates a key for the current base using credentials found
 in the global config file. The key is then inserted into ./.cfg/key.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		key.Create(team, base, name, canRead, canWrite, expiry)
	},
}

func init() {
	createCmd.Flags().StringVarP(&team, "team", "t", "", "name of team")
	createCmd.Flags().StringVarP(&base, "base", "b", "", "name of base")
	createCmd.Flags().StringVarP(&name, "name", "n", "", "name of key")
	createCmd.Flags().BoolVarP(&canRead, "can-read", "", false, "specify key read access")
	createCmd.Flags().BoolVarP(&canWrite, "can-write", "", false, "specify key write access")
	createCmd.Flags().IntVarP(&expiry, "expiry", "e", 0, "expiry timestamp (seconds since epoch; 0 = never)")

	keyCmd.AddCommand(createCmd)
}
