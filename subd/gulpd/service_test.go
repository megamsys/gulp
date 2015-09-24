/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
 
package gulpd

import (
	"gopkg.in/check.v1"
	"github.com/megamsys/gulp/meta"
)

type S struct {
	service *Service
}

var _ = check.Suite(&S{})

func (s *S) TestCreateService(c *check.C) {

    srv := NewService(&meta.Config{}, &Config{
		CatID: "ASM000",
	})	
	s.service = srv
	c.Assert(srv, check.NotNil)
}
