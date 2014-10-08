package main

import (
	"bytes"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/cmd/testing"
//	"io/ioutil"
	"gopkg.in/check.v1"
	"net/http"
//	"strings"
)

func (s *S) TestAppStartInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "startapp",
		Usage:   "startapp <appname> <lifecycle_when>",
		Desc:    "starts the installed app.",
		MinArgs: 1,
	}
	c.Assert((&AppStart{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestAppStart(c *check.C) {
	var stdout, stderr bytes.Buffer
//	result := `{"status":"success", "repository_url":"git@github.com/indykish:nilavu.git"}`
	expected := `App "ble.megam.co" is being started!
Use appreqs list to check the status of the app.` + "\n"
	context := cmd.Context{
		Args:   []string{"ble.megam.co", "rails"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
/*	trans := testing.ConditionalTransport{
		Transport: testing.Transport{Message: result, Status: http.StatusOK},
		CondFunc: func(req *http.Request) bool {
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			c.Assert(err, check.IsNil)
			c.Assert(string(body), check.Equals, `{"name":"ble","platform":"django"}`)
			return req.Method == "POST" && req.URL.Path == "/apps"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	*/
	command := AppStart{}
	trans := testing.Transport{Message: "success", Status: http.StatusOK}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	err := command.Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, expected)
}
