package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/indykish/gulp/fs"
	"github.com/indykish/gulp/fs/testing"
	"io"
	"launchpad.net/gnuflag"
	"launchpad.net/gocheck"
	"os"
)

type recordingExiter int

func (e *recordingExiter) Exit(code int) {
	*e = recordingExiter(code)
}

func (e recordingExiter) value() int {
	return int(e)
}

type TestCommand struct{}

func (c *TestCommand) Info() *Info {
	return &Info{
		Name:  "foo",
		Desc:  "Foo do anything or nothing.",
		Usage: "foo",
	}
}

func (c *TestCommand) Run(context *Context) error {
	io.WriteString(context.Stdout, "Running TestCommand")
	return nil
}

type ErrorCommand struct {
	msg string
}

func (c *ErrorCommand) Info() *Info {
	return &Info{Name: "error"}
}

func (c *ErrorCommand) Run(context *Context) error {
	return errors.New(c.msg)
}

func (s *S) TestRegister(c *gocheck.C) {
	manager.Register(&TestCommand{})
	badCall := func() { manager.Register(&TestCommand{}) }
	c.Assert(badCall, gocheck.PanicMatches, "command already registered: foo")
}

func (s *S) TestManagerRunShouldWriteErrorsOnStderr(c *gocheck.C) {
	manager.Register(&ErrorCommand{msg: "You are wrong\n"})
	manager.Run([]string{"error"})
	c.Assert(manager.stderr.(*bytes.Buffer).String(), gocheck.Equals, "Error: You are wrong\n")
}


func (s *S) TestRun(c *gocheck.C) {
	manager.Register(&TestCommand{})
	manager.Run([]string{"foo"})
	c.Assert(manager.stdout.(*bytes.Buffer).String(), gocheck.Equals, "Running TestCommand")
}


func (s *S) TestFileSystem(c *gocheck.C) {
	fsystem = &testing.RecordingFs{}
	c.Assert(filesystem(), gocheck.DeepEquals, fsystem)
	fsystem = nil
	c.Assert(filesystem(), gocheck.DeepEquals, fs.OsFs{})
}

