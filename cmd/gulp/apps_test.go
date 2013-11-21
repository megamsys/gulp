package main

import (
	"bytes"
	"github.com/indykish/gulp/cmd"
	"github.com/indykish/gulp/cmd/testing"
	"io/ioutil"
	"launchpad.net/gocheck"
	"net/http"
	"strings"
)

func (s *S) TestAppCreateInfo(c *gocheck.C) {
	expected := &cmd.Info{
		Name:    "app-create",
		Usage:   "app-create <appname> <platform>",
		Desc:    "create a new app.",
		MinArgs: 2,
	}
	c.Assert((&AppCreate{}).Info(), gocheck.DeepEquals, expected)
}

func (s *S) TestAppCreate(c *gocheck.C) {
	var stdout, stderr bytes.Buffer
	result := `{"status":"success", "repository_url":"git@github.com/indykish:nilavu.git"}`
	expected := `App "ble" is being created!
Use app-info to check the status of the app and its units.
Your repository for "ble" project is "git@github.com/indykish:ble.git"` + "\n"
	context := cmd.Context{
		Args:   []string{"ble", "django"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
/*	trans := testing.ConditionalTransport{
		Transport: testing.Transport{Message: result, Status: http.StatusOK},
		CondFunc: func(req *http.Request) bool {
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			c.Assert(err, gocheck.IsNil)
			c.Assert(string(body), gocheck.Equals, `{"name":"ble","platform":"django"}`)
			return req.Method == "POST" && req.URL.Path == "/apps"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
*/	command := AppCreate{}
	err := command.Run(&context)
	c.Assert(err, gocheck.IsNil)
	c.Assert(stdout.String(), gocheck.Equals, expected)
}
