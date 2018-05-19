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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "remove the remote with the given name",
	Long: `Removes the remote with the given name.

Example: cfg remote rm origin`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remote.Remove(args[0])
	},
}

func init() {
	remoteCmd.AddCommand(rmCmd)
}
