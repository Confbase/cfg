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

	"github.com/Confbase/cfg/lib/snap"
)

var snapCheckoutCmd = &cobra.Command{
	Use:   "checkout <snapshot-name>",
	Short: "Checkout a snapshot",
	Long: `Checks out a snapshot.

If the base does not use git, this command is disabled. A post-checkout hook has
already been added to the surrounding git repository to checkout the snapshot
automatically. If there is no surrounding git repository, nothing happens and
a warning is displayed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snap.Checkout(args[0])
	},
}

func init() {
	snapCmd.AddCommand(snapCheckoutCmd)
}
