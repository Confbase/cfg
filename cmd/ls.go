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

	"github.com/Confbase/cfg/ls"
)

var lsCfg ls.Config
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List the contents of the base",
	Long: `Lists the contents of the base.

Use -t, -i, -s to only list templates, instances, and singletons, respectively.

If stdout is a tty, then the contents are listed in a human-readable format.

If stdout is not a tty, or if the --no-tty flag is used, then the contents
are listed in the following format:

templates
<template-name><tab character><template-file-path>
<template-name><tab character><template-file-path>
...
instances
<instance-file-path><tab character><template-name>,<template-name>,...
<instance-file-path><tab character><template-name>,<template-name>,...
....
singletons
<singleton-filepath>
<singleton-filepath>
...
`,
	Run: func(cmd *cobra.Command, args []string) {
		ls.Ls(&lsCfg)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&lsCfg.NoTty, "no-tty", "", false, "do not print in human-readable format")
	lsCmd.Flags().BoolVarP(&lsCfg.NoColors, "no-colors", "", false, "do not colorize text")
	lsCmd.Flags().BoolVarP(&lsCfg.DoLsTempls, "templates", "t", false, "only list templates")
	lsCmd.Flags().BoolVarP(&lsCfg.DoLsInsts, "instances", "i", false, "only list instances")
	lsCmd.Flags().BoolVarP(&lsCfg.DoLsSingles, "singletons", "s", false, "only list singletons")
	lsCmd.Flags().BoolVarP(&lsCfg.DoLsUntracked, "untracked", "u", false, "only list untracked")
}
