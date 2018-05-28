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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Confbase/cfg/clone"
)

var cloneCfg clone.Config
var cloneCmd = &cobra.Command{
	Use:   "clone <base> [<target directory>]",
	Short: "Clone a base",
	Long: `Clones a base. Can be used with or without --no-git.

If a directory is specified, the base will be cloned into that directory.

Otherwise, the base will be cloned into ./<base name>.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 1:
			cloneCfg.RemoteURL = args[0]
		case 2:
			cloneCfg.RemoteURL = args[0]
			cloneCfg.TargetDir = args[1]
		default:
			fmt.Fprintf(os.Stderr, "usage: cfg clone <base> [<target directory>]\n")
			os.Exit(1)
		}
		clone.Clone(&cloneCfg)
	},
}

func init() {
	RootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().BoolVarP(&cloneCfg.NoGit, "no-git", "", false, "do not use git in cloned base")
	cloneCmd.Flags().StringVarP(&cloneCfg.DefaultSnap, "default-snap", "s", "master", "specifies snapshot to checkout first")
}
