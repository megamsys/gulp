package policies

import (
	"launchpad.net/gocheck"
)

func (s *S) TestRegisterPolicies(c *gocheck.C) {
	provider, err := GetPolicy("abc")
	c.Assert(err, gocheck.ErrorMatches, "policy \"abc\" not registered")
	c.Assert(provider, gocheck.IsNil)
	policy := TestIaaS{}
	RegisterPolicy("abc", policy)
	p, err = getIaasProvider("abc")
	c.Assert(err, gocheck.IsNil)
	c.Assert(p, gocheck.Equals, policy)
}

