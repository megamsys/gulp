
package handlers

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestHandler(c *check.C) {
	error := Handler("CAT1248571413820997632")
	c.Assert(error, check.IsNil)
}
