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
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cfg",
	Short: "The official Confbase CLI",
	Long: `cfg is the official CLI for Confbase (https://confbase.com).

It is a wrapper around git, but provides additionaly functionality in the
form of templates, instances, singletons, and numerous operations on these
objects. Furthermore, a raw copy of every file in the latest commit is made
available via 'cfg fetch'.

cfg intends to be brutalist. All work is done on one branch. When the remote
branch gets ahead of the local branch, 'cfg trampoline' or 'cfg pull --hard'
are the only two options to fix the predicament.

If a more traditional workflow is desired, consider moving the configuration
files back into their codebase's repository and simply using git as before.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "global-config", "", "global config file (default is $HOME/.cfg.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".cfg")  // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	viper.SetDefault("email", "")
	viper.SetDefault("password", "")
	viper.SetDefault("entryPoint", "https://confbase.com")

	viper.ReadInConfig()
}
