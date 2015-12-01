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

package gitlab

import (
	//	"bytes"
	//	"net/http"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&GitlabSuite{})

type GitlabSuite struct {
}

/*func (s *GitlabSuite) SetUpSuite(c *check.C) {
	var err = error.New("testing")
	c.Assert(err, check.IsNil)
}

func (s *GitlabSuite) TestClone(c *check.C) {
	var err = error.New("testing")
	c.Assert(err, check.IsNil)
}*/
