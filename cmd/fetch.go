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

	"github.com/Confbase/cfg/fetch"
)

var snapshotToFetch, fetchEmail, fetchPassword string
var fetchNoGlobal bool

var fetchCmd = &cobra.Command{
	Use:   "fetch [targets]",
	Short: "Fetch the latest copy of a file",
	Long: `Fetches the latest copy of a file without
downloading the entire base.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: consider "pipelining" these requests so we only make one request
		// TODO: consider moving args into a config struct, since there are so many
		for _, target := range args {
			fetch.Fetch(snapshotToFetch, target, fetchEmail, fetchPassword, fetchNoGlobal)
		}
	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&snapshotToFetch, "snapshot", "s", "master", "specifies snapshot")
	fetchCmd.Flags().StringVarP(&fetchEmail, "email", "", "", "specifies email")
	fetchCmd.Flags().StringVarP(&fetchPassword, "password", "", "", "specifies password")
	fetchCmd.Flags().BoolVarP(&fetchNoGlobal, "no-global-config", "n", false, "do not use global config file")
}
