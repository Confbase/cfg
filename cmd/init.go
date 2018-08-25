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
	initpkg "github.com/Confbase/cfg/init"

	"github.com/spf13/cobra"
)

var initCfg initpkg.Config

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new base in the current directory",
	Long: `Initializes a new base, using credentials specified in the
global config file.

By default, cfg creates a new git repository and controls it
when certain cfg commands are run, like 'cfg push'.

If you wish to use cfg inside an existing git repository, then
pass the --no-git flag to 'cfg init'. This will stop cfg from
creating a git repository and it will stop cfg from issuing git
commands in this base.

It is possible to force the creation of a git-backed base with
'cfg init --force'.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			initCfg.Dest = args[0]
		}
		initpkg.MustInit(&initCfg)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&initCfg.OverwriteGitIgnore, "overwrite-gitignore", "", false, "overwrite .gitignore")
	initCmd.Flags().BoolVarP(&initCfg.AppendGitIgnore, "append-to-gitignore", "", false, "append to .gitignore")
	initCmd.Flags().BoolVarP(&initCfg.NoGit, "no-git", "", false, "do not create or control a git repository")
	initCmd.Flags().BoolVarP(&initCfg.NoModGitIgnore, "no-modify-gitignore", "", false, "do not modify .gitignore")
	initCmd.Flags().BoolVarP(&initCfg.Force, "force", "f", false, "force the creation of a base")
}
