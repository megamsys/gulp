package scm

import (
	"github.com/globocom/config"
	"launchpad.net/gocheck"
)


func (s *S) TestGetPath(c *gocheck.C) {
	path, err := GetPath()
	c.Assert(err, gocheck.IsNil)
	expected := "/home/application/current"
	c.Assert(path, gocheck.Equals, expected)
}

func (s *S) TestGetServerUri(c *gocheck.C) {
	server, err := config.GetString("scm:api_server")
	c.Assert(err, gocheck.IsNil)
	uri := ServerURL()
	c.Assert(uri, gocheck.Equals, server)
}

func (s *S) TestGetServerUriWithoutSetting(c *gocheck.C) {
	old, _ := config.Get("scm:api_server")
	defer config.Set("scm:api_server", old)
	config.Unset("scm:api_server")
	defer func() {
		r := recover()
		c.Assert(r, gocheck.NotNil)
	}()
	ServerURL()
}