
package control

import (
	"testing"
"github.com/megamsys/gulp"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})



func (s *S) TestControlHandler(c *check.C) {
  req := &handlers.Request{Id: "CAT1248571413820997632"}
	error := ControlHandler(req)
	c.Assert(error, check.IsNil)
}
