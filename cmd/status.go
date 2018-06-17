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

	"github.com/Confbase/cfg/status"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the current base",
	Long:  `Shows the status of the current base.`,
	Run: func(cmd *cobra.Command, args []string) {
		status.MustShowStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
