package cmdrunner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func nilOrFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func RunOrFatal(cmd *exec.Cmd) {
	cmdStderr, err := cmd.StderrPipe()
	nilOrFatal(err)
	cmdStdout, err := cmd.StdoutPipe()
	nilOrFatal(err)
	err = cmd.Start()
	nilOrFatal(err)

	go func() {
		_, err := io.Copy(os.Stdout, cmdStdout)
		nilOrFatal(err)
	}()
	go func() {
		_, err := io.Copy(os.Stderr, cmdStderr)
		nilOrFatal(err)
	}()

	err = cmd.Wait()
	nilOrFatal(err)
}
