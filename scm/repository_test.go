package scm

import (
	"github.com/globocom/config"
	"launchpad.net/gocheck"
)


func (s *S) TestGetPath(c *gocheck.C) {
	path, err := GetPath()
	c.Assert(err, gocheck.IsNil)
	expected := "/var/www/projects/aryabhata/current"
	c.Assert(path, gocheck.Equals, expected)
}

func (s *S) TestGetRemotePath(c *gocheck.C) {
	path, err := GetRemotePath()
	c.Assert(err, gocheck.IsNil)
	expected := "https://github.com/indykish/aryabhata.git"
	c.Assert(path, gocheck.Equals, expected)
}

func (s *S) TestProject(c *gocheck.C) {
	path, err := Project()
	c.Assert(err, gocheck.IsNil)
	expected := "aryabhata"
	c.Assert(path, gocheck.Equals, expected)
}

func (s *S) TestBuilder(c *gocheck.C) {
	path, err := Builder()
	c.Assert(err, gocheck.IsNil)
	expected := "megam_builder_ruby"
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