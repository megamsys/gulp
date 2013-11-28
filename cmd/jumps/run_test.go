package jumps

import (
	"bytes"
	"github.com/indykish/gulp/cmd"
	"launchpad.net/gocheck"
)



func (s *S) TestGulpStartInfo(c *gocheck.C) {
	desc := `run a command in all instances of the app, and prints the output.

If you use the '--once' flag tsuru will run the command only in one unit.

If you don't provide the app name, tsuru will try to guess it.
`
	expected := &cmd.Info{
		Name:    "run",
		Usage:   `run <command> [commandarg1] [commandarg2] ... [commandargn] [--app appname] [--once]`,
		Desc:    desc,
		MinArgs: 1,
	}
	command := AppRun{}
	c.Assert(command.Info(), gocheck.DeepEquals, expected)
}


func (s *S) TestGulpStopInfo(c *gocheck.C) {
	desc := `run a command in all instances of the app, and prints the output.

If you use the '--once' flag tsuru will run the command only in one unit.

If you don't provide the app name, tsuru will try to guess it.
`
	expected := &cmd.Info{
		Name:    "run",
		Usage:   `run <command> [commandarg1] [commandarg2] ... [commandargn] [--app appname] [--once]`,
		Desc:    desc,
		MinArgs: 1,
	}
	command := AppRun{}
	c.Assert(command.Info(), gocheck.DeepEquals, expected)
}

