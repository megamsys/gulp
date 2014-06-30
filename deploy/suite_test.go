package deploy

import (
	"github.com/tsuru/config"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	config.Set("git:unit-repo", "test/dir")
	config.Set("git:ro-host", "api.megam.co")
}

func (s *S) TearDownSuite(c *check.C) {
	config.Unset("git:unit-repo")
	config.Unset("git:host")
}
