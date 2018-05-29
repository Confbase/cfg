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
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var doCheckoutB bool
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "alias for 'cfg snap checkout'",
	Long:  `Alias for 'cfg snap checkout'. See 'cfg snap checkout --help' for more info.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cobraCmd *cobra.Command, args []string) {
		cmdArgs := make([]string, 2)
		cmdArgs[0] = "snap"
		if doCheckoutB {
			cmdArgs[1] = "new"
		} else {
			cmdArgs[1] = "checkout"
		}
		cmdArgs = append(cmdArgs, args...)
		cmd := exec.Command("cfg", cmdArgs...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		go func() {
			_, err := io.Copy(os.Stdout, stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		}()
		go func() {
			_, err := io.Copy(os.Stderr, stderr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		}()

		if err := cmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				// The program has exited with an exit code != 0
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
	checkoutCmd.Flags().BoolVarP(&doCheckoutB, "create-snap", "b", false, "alias for 'cfg snap new'")
}
