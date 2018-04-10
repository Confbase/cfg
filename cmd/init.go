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
	initpkg "github.com/Confbase/cfg/lib/init"

	"github.com/spf13/cobra"
)

var initAppend, initOverwrite bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new base in the current directory",
	Long: `Initializes a new base, using credentials
specified in the global config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		initpkg.Init(initAppend, initOverwrite)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&initOverwrite, "overwrite-gitignore", "", false, "overwrite .gitignore")
	initCmd.Flags().BoolVarP(&initAppend, "append-to-gitignore", "", false, "append to .gitignore")
}
