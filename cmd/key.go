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

import "github.com/spf13/cobra"

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Per-base key management",
	Long: `This command helps manage access keys at a per-base level.

Interacting with Confbase servers requires authentication with
access keys. There are two ways to acquire and use access keys.

1. cfg key create

This generates a new access key for this base and inserts it into
./.cfg.json, using the credentials found in ~/.cfg.yaml.

2. cfg key insert

After generating a key through the web interface, the key can be
inserted into ./.cfg.json with this command.`,
}

func init() {
	RootCmd.AddCommand(keyCmd)
}
