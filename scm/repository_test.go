package scm

import (
	"github.com/tsuru/config"
	"gopkg.in/check.v1"
)

func (s *S) TestGetPath(c *check.C) {
	path, err := GetPath()
	c.Assert(err, check.IsNil)
	expected := "/var/www/projects/aryabhata/current"
	c.Assert(path, check.Equals, expected)
}

func (s *S) TestGetRemotePath(c *check.C) {
	path, err := GetRemotePath()
<<<<<<< HEAD
	c.Assert(err, gocheck.IsNil)
	expected := "https://github.com/megamsys/aryabhata.git"
	c.Assert(path, gocheck.Equals, expected)
=======
	c.Assert(err, check.IsNil)
	expected := "https://github.com/indykish/aryabhata.git"
	c.Assert(path, check.Equals, expected)
>>>>>>> origin/master
}

func (s *S) TestProject(c *check.C) {
	path, err := Project()
	c.Assert(err, check.IsNil)
	expected := "aryabhata"
	c.Assert(path, check.Equals, expected)
}

func (s *S) TestBuilder(c *check.C) {
	path, err := Builder()
	c.Assert(err, check.IsNil)
	expected := "megam_builder_ruby"
	c.Assert(path, check.Equals, expected)
}

func (s *S) TestGetServerUri(c *check.C) {
	server, err := config.GetString("scm:api_server")
	c.Assert(err, check.IsNil)
	uri := ServerURL()
	c.Assert(uri, check.Equals, server)
}

func (s *S) TestGetServerUriWithoutSetting(c *check.C) {
	old, _ := config.Get("scm:api_server")
	defer config.Set("scm:api_server", old)
	config.Unset("scm:api_server")
	defer func() {
		r := recover()
		c.Assert(r, check.NotNil)
	}()
	ServerURL()
}
