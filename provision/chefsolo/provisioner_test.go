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

package chefsolo

import (
	//	"os"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})

/*
func (s *S) TestPrepareFiles(c *check.C) {
		os.MkdirAll("/tmp/chef-solo/sandbox", 0755)

		p := chefsoloProvisioner{
			SandboxPath: "/tmp/chef-solo/sandbox",
			RootPath:    "/tmp/chef-solo",
		}
		c.Assert(p.PrepareFiles(), check.IsNil)
}

func (s *S) TestPrepareFiles_CustomJSON(c *check.C) {
		os.MkdirAll(".chef-solo/sandbox", 0755)

		p := chefsoloProvisioner{
			Attributes:  `{"foo": "bar"}`,
			SandboxPath: ".chef-solo/sandbox",
			RootPath:    "/tmp/chef-solo",
		}
	c.Assert(p.PrepareFiles(), check.IsNil)

}*/
