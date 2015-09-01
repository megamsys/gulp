
package handler

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
	error := Handler("CAT1242820117956526080")
	c.Assert(error, check.IsNil)
}
