package scm

import (
	"github.com/globocom/config"
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	config.Set("scm:repo", "git:github.com/indykish/aryabhata.git")
	config.Set("scm:api_server", "api.github.com/v1")
	config.Set("scm:local_repo", "/home/application/current")
}