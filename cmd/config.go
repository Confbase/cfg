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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Confbase/cfg/config"
)

var configCmd = &cobra.Command{
	Use:   "config [<key> <value>]",
	Short: "Configure cfg variables per base",
	Long: `Configures cfg variables per base.

Example---set email address for this base:

    $ cfg config email thomas@fr.tld

Example---set base name to "myproject":

    $ cfg config base myproject`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args)%2 != 0 {
			fmt.Fprintf(os.Stderr, "error: 'cfg config' expects key-value pairs\n")
			os.Exit(1)
		}
		config.Config(args)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
