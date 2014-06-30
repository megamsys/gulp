
package heal

import (
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})
/*
func (s *S) TestRegisterAndGetHealer(c *check.C) {
	var h Healer
	Register("my-provisioner", "my-healer", h)
	got, err := Get("my-provisioner", "my-healer")
	c.Assert(err, check.IsNil)
	c.Assert(got, check.DeepEquals, h)
	_, err = Get("my-provisioner", "unknown-healer")
	c.Assert(err, check.ErrorMatches, `Unknown healer "unknown-healer" for provisioner "my-provisioner".`)
}

func (s *S) TestGetWithAbsentProvisioner(c *check.C) {
	var h Healer
	Register("provisioner", "healer1", h)
	h, err := Get("otherprovisioner", "healer1")
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, `Unknown healer "healer1" for provisioner "otherprovisioner".`)
	c.Assert(h, check.IsNil)
}

func (s *S) TestAllReturnsAllByCurrentProvisioner(c *check.C) {
	var h Healer
	Register("provisioner", "healer1", h)
	Register("provisioner", "healer2", h)
	healers := All("provisioner")
	expected := map[string]Healer{
		"healer1": h,
		"healer2": h,
	}
	c.Assert(healers, check.DeepEquals, expected)
}
*/