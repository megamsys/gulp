package main

import (
	"github.com/indykish/gulp/cmd"
	"launchpad.net/gocheck"
)

func (s *S) TestCommandsFromBaseManagerAreRegistered(c *gocheck.C) {
	baseManager := cmd.BuildBaseManager("megam", version, header)
	manager := buildManager("megam")
	for name, instance := range baseManager.Commands {
		command, ok := manager.Commands[name]
		c.Assert(ok, gocheck.Equals, true)
		c.Assert(command, gocheck.FitsTypeOf, instance)
	}
}

func (s *S) TestAppCreateIsRegistered(c *gocheck.C) {
	manager := buildManager("megam")
	create, ok := manager.Commands["app-create"]
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(create, gocheck.FitsTypeOf, AppCreate{})
}

func (s *S) TestAppRemoveIsRegistered(c *gocheck.C) {
	manager := buildManager("megam")
	remove, ok := manager.Commands["app-remove"]
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(remove, gocheck.FitsTypeOf, &AppRemove{})
}

func (s *S) TestAppListIsRegistered(c *gocheck.C) {
	manager := buildManager("megam")
	list, ok := manager.Commands["app-list"]
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(list, gocheck.FitsTypeOf, megam.AppList{})
}