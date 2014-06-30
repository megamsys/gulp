package scm

import (
	"github.com/tsuru/config"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	config.Set("scm:api_server", "api.github.com/v1")
	config.Set("scm:local_repo", "/var/www/projects/aryabhata/current")
	config.Set("scm:remote_repo", "https://github.com/indykish/aryabhata.git")
	config.Set("scm:project", "aryabhata")
	config.Set("scm:builder", "megam_builder_ruby")

}
