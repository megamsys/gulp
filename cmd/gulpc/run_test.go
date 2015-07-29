package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
)

func (s *S) TestGulpcStartInfo(c *check.C) {
	desc := `starts the gulpc server.

If you use the '--dry' flag gulpc server will do a dry run(parse conf/jsons) and exit.

`

	expected := &cmd.Info{
		Name:    "start",
		Usage:   `start [--dry] [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := GulpcStart{}
	c.Assert(command.Info(), check.DeepEquals, expected)
}

func (s *S) TestGulpcStopInfo(c *check.C) {
	desc := `stops the gulpc daemon, and shutsdown the queue.

If you use the '--bark' flag gulpc will notify daemon status.

`
	expected := &cmd.Info{
		Name:    "stop",
		Usage:   `stop [--bark]`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := GulpcStop{}
	c.Assert(command.Info(), check.DeepEquals, expected)
}


