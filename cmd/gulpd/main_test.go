package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})



func (s *S) TestCommandsFromBaseManagerAreRegistered(c *check.C) {
	baseManager := cmd.BuildBaseManager("megam", version, header)
	manager := buildManager("megam")
	for name, instance := range baseManager.Commands {
		command, ok := manager.Commands[name]
		c.Assert(ok, check.Equals, true)
		c.Assert(command, check.FitsTypeOf, instance)
	}
}


