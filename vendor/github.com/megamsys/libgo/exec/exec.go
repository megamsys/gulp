// Package exec provides a interface to run external commans as an
// abstraction layer.
package exec

import (
	"io"
	"os/exec"
)

type Executor interface {
	// Execute executes the specified command.
	Execute(cmd string, args []string, stdin io.Reader, stdout, stderr io.Writer) error
}

type OsExecutor struct{}

func (OsExecutor) Execute(cmd string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	c := exec.Command(cmd, args...)
	c.Stdin = stdin
	c.Stdout = stdout
	c.Stderr = stderr
	return c.Run()
}
