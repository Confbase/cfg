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

	"github.com/Confbase/cfg/remote"
)

// lsCmd represents the ls command
var remoteLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list remotes",
	Long:  `Lists the remotes in the local keyfile.`,
	Run: func(cmd *cobra.Command, args []string) {
		remote.Ls()
	},
}

func init() {
	remoteCmd.AddCommand(remoteLsCmd)
}
