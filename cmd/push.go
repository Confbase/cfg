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

	"github.com/Confbase/cfg/push"
)

var pushCfg push.Config

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to host server",
	Long:  `Pushes changes to host server.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 1:
			pushCfg.Remote = args[0]
		case 2:
			pushCfg.Remote = args[0]
			pushCfg.Snapshot = args[1]
		}
		push.Push(pushCfg)
	},
}

func init() {
	pushCmd.Flags().BoolVarP(&pushCfg.IsForce, "force", "f", false, "force push")
	rootCmd.AddCommand(pushCmd)
}
