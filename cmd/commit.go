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

	"github.com/Confbase/cfg/commit"
)

var commitCfg commit.Config
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes made",
	Long: `Commit changes made to the contents of the base.

This command is not available in --no-git bases.

If no files are in the staging area, all tracked files staged before the commit.

If files are already in the staging area, only those are included in the commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		commit.MustCommitOrRevert(&commitCfg)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&commitCfg.Message, "message", "m", "", "commit message")
}
