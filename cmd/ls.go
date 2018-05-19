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

var lsNoTty, lsNoColors bool

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List the contents of the base",
	Long:  `Lists the contents of the base.

If stdout is a tty, then the contents are listed in a human-readable format.

If stdout is not a tty, or if the --no-tty flag is used, then the contents
are listed in the following format:

templates
<template-name>,<template-file-type>,<template-file-path>
<template-name>,<template-file-type>,<template-file-path>
...
instances
<template-name>,<instance-file-path>
<template-name>,<instance-file-path>
....
singletons
<singleton-filepath>
<singleton-filepath>
...

That is, the string literal "templates", followed by a newline character and
the templates, one per line and in CSV format; followed by the string literal
"instances", followed by a newline character and the instances, one per line
and in CSV format; followed by the string literal "singletons\n", followed by a
newline character and the singletons, one per line.
`,
	Run: func(cmd *cobra.Command, args []string) {
		ls.Ls(lsNoTty, lsNoColors)
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&lsNoTty, "no-tty", "", false, "do not print in human-readable format")
	lsCmd.Flags().BoolVarP(&lsNoColors, "no-colors", "", false, "do not colorize text")
}
