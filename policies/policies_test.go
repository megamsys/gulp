package policies

import (
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}


var _ = check.Suite(&S{})


func (s *S) TestRegisterPolicies(c *check.C) {
	provider, _ := GetPolicy("abc")
	c.Assert(provider, check.IsNil)
}

