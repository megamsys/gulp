package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
)

func (s *S) TestGulpdStartInfo(c *check.C) {
	desc := `starts the gulpd daemon, and connects to queue.

If you use the '--dry' flag gulpd will do a dry run(parse conf/jsons) and exit.

`

	expected := &cmd.Info{
		Name:    "start",
		Usage:   `start [--dry] [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := GulpdStart{}
	c.Assert(command.Info(), check.DeepEquals, expected)
}


func (s *S) TestGulpdStopInfo(c *check.C) {
	desc := `stops the gulpd daemon, and shutsdown the queue.

If you use the '--bark' flag gulpd will notify daemon status.

`
	expected := &cmd.Info{
		Name:    "stop",
		Usage:   `stop [--bark]`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := GulpdStop{}
	c.Assert(command.Info(), check.DeepEquals, expected)
}
