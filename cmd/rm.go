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

	"github.com/Confbase/cfg/rm"
)

// Configuration for the running of rm
var rmCfg rm.Config

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a file and unmark it",
	Long: `Removes a file and unmarks it.
	
	By default, this will only operate on a single file, however, it 
	is possible to recursively call 'cfg rm' on a directory and all of
	its children through using the -r flag
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rmCfg.ToRemove = args
		rm.MustRm(&rmCfg)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&rmCfg.Recursive, "recursive", "r", false, "recursively call rm on all children of target")
}
