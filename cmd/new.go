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

	"github.com/Confbase/cfg/new"
)

var newCfg new.Config
var newCmd = &cobra.Command{
	Use:   "new <template-name> [target-path]",
	Short: "Create a new instance of a template",
	Long: `Creates a new instance of a template.

If no targets are specified, the instance is written to stdout.

If targets are specified, a new instance is initialized at each target path.

Target paths starting with the string "-" or "--" ARE IGNORED for the
following reason:

This command determines the location of the given template's schema and then
runs 'schema init -s <schema-path> [target-path...]'. This means that it is
possible to use any of the flags 'schema path' supports by prefixing them
with the target "--". See the next two sections for examples.

SPECIFYING FORMAT

$ cfg new myTempl myInstance.toml -- --toml

To see all supported formats, run 'schema init -h'.

INITIALIZING WITH RANDOM VALUES

$ cfg new myTempl myInstance.json -- --random`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			newCfg.Targets = args[1:]
		}
		newCfg.TemplName = args[0]
		new.MustNew(&newCfg)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&newCfg.DoUseDefaults, "use-default-values", "d", false, "intialize fields with default values")
}
