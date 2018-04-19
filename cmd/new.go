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

var newCmd = &cobra.Command{
	Use:   "new <snapshot-name>",
	Short: "Associate a copy of the current files with a given name",
	Long: `Creates a snapshot of the current files and associates it with a name.

These files are stored in an in-memory cache on Confbase servers. They can be
accessed on any computer via

$ cfg fetch myteam.confbase.com:mybase/myfile.json --snapshot my-snapshot

or, if cfg is not on the machine, simply via HTTPS

$ curl -u 'key:xxxxxxxxx' myteam.confbase.com/mybase/raw/my-snapshot/myfile.json
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snap.New(args[0])
	},
}

func init() {
	snapCmd.AddCommand(newCmd)
}
