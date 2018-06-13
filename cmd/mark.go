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

	"github.com/Confbase/cfg/mark"
)

var markCfg mark.Config

var markCmd = &cobra.Command{
	Use:   "mark <file> <template-name>",
	Short: "Mark a file as a template, instance, or singleton",
	Long: `Mark a file as a template, instance, or singleton.

--template,    -t specifies the file is a template
--instance-of, -i specifies the file is an instance of a template
--singleton,   -s specifies the file is a singleton
--unmark,      -u unmarks the file

Using more than one of the above four flags simultaneously is
undefined.

New instances of template files can be created with "cfg new [template-file]".`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		markCfg.Targets = args
		mark.Mark(&markCfg)
	},
}

func init() {
	rootCmd.AddCommand(markCmd)
	markCmd.Flags().BoolVarP(&markCfg.Force, "force", "", false, "overwrite template if it already exists")
	markCmd.Flags().BoolVarP(&markCfg.UnMark, "unmark", "u", false, "unmark file")
	markCmd.Flags().StringVarP(&markCfg.InstanceOf, "instance-of", "i", "", "specifies template which target is instance of")
	markCmd.Flags().StringVarP(&markCfg.Template, "template", "t", "", "specifies template name associated with target")
	markCmd.Flags().BoolVarP(&markCfg.Singleton, "singleton", "s", false, "track target as singleton")
}
