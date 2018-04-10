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

	"github.com/confbase/cfg/lib/mark"
)

var markForce bool

var markCmd = &cobra.Command{
	Use:   "mark <file> <template-name>",
	Short: "Mark a file as a template file",
	Long: `Marks a file as a template file.

New instances of template files can be created with "cfg new [template-file]".

See related "cfg tag" command.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mark.Mark(args[0], args[1], markForce)
	},
}

func init() {
	RootCmd.AddCommand(markCmd)
	markCmd.Flags().BoolVarP(&markForce, "force", "", false, "overwrite template if it already exists")
}
