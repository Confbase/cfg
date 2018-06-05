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

func PipeFrom(cmd *exec.Cmd, wOut io.Writer, wErr io.Writer) error {
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		if wOut != nil {
			io.Copy(wOut, cmdStdout)
		}
	}()
	go func() {
		if wErr != nil {
			io.Copy(wErr, cmdStderr)
		}
	}()
	return cmd.Wait()
}

func PipeThrough(cmd *exec.Cmd, rIn io.Reader, wOut io.Writer, wErr io.Writer) error {
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		if rIn != nil {
			io.Copy(cmdStdin, rIn)
		}
		cmdStdin.Close()
	}()
	go func() {
		if wOut != nil {
			io.Copy(wOut, cmdStdout)
		}
	}()
	go func() {
		if wErr != nil {
			io.Copy(wErr, cmdStderr)
		}
	}()
	return cmd.Wait()
}
