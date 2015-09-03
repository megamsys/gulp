package gulpd

import (
	"github.com/megamsys/gulp/meta"
	"gopkg.in/check.v1"
)

type S struct {
	service *activity.Service
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	srv, err := &NewService(activity.Config{
		BindAddress: "127.0.0.1:0",
	})
	s.service = srv
	c.Assert(err, check.IsNil)
}
