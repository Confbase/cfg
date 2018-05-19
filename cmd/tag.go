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

	"github.com/Confbase/cfg/tag"
)

var tagCmd = &cobra.Command{
	Use:   "tag <file> <template-name>",
	Short: "Tag a file as an instance of a template",
	Long: `Tags a file as an instance of a certain template.

Future releases will support group operations on all instances of 
a certain template. Files can also be sorted by tag on the Confbase
web interface.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tag.Tag(args[0], args[1])
	},
}

func init() {
	RootCmd.AddCommand(tagCmd)
}
