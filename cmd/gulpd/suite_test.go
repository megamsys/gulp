package main

import (
	"bytes"
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
	"os"
	"os/exec"
	"testing"
)

type S struct {
	recover []string
}

func (s *S) SetUpSuite(c *check.C) {
	targetFile := os.Getenv("HOME") + "/.megam"
	_, err := os.Stat(targetFile)
	if err == nil {
		old := targetFile + ".old"
		s.recover = []string{"mv", old, targetFile}
		exec.Command("mv", targetFile, old).Run()
	} else {
		s.recover = []string{"rm", targetFile}
	}
	f, err := os.Create(targetFile)
	c.Assert(err, check.IsNil)
	f.Write([]byte("http://localhost"))
	f.Close()
}

func (s *S) TearDownSuite(c *check.C) {
	exec.Command(s.recover[0], s.recover[1:]...).Run()
}

var _ = check.Suite(&S{})
var manager *cmd.Manager

func Test(t *testing.T) { check.TestingT(t) }

func (s *S) SetUpTest(c *check.C) {
	var stdout, stderr bytes.Buffer
	manager = cmd.NewManager("gulpd", version, header, &stdout, &stderr, os.Stdin)
}
