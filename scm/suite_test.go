package scm

import (
	"github.com/tsuru/config"
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	config.Set("scm:api_server", "api.github.com/v1")
	config.Set("scm:local_repo", "/var/www/projects/aryabhata/current")
	config.Set("scm:remote_repo", "https://github.com/megamsys/aryabhata.git")
	config.Set("scm:project", "aryabhata")
	config.Set("scm:builder", "megam_builder_ruby")

}
