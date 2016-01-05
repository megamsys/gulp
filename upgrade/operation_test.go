/*
** Copyright [2013-2016] [Megam Systems]
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

package upgrade

/*import (
	"gopkg.in/check.v1"
  "testing"
)
type S struct{}

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&S{})

func (s *S) TestRegister(c *check.C) {
	mngr := nopManager{}
	Register("bind", mngr)
	defer func() {
		delete(managers, "bind")
	}()
	c.Assert(managers["bind"], check.Equals, mngr)
}

/*func (s *S) TestRegisterOnNilMap(c *check.C) {
	oldManagers := managers
	managers = nil
	defer func() {
		managers = oldManagers
	}()
	mngr := nopManager{}
	Register("nope", mngr)
	c.Assert(managers["nope"], check.Equals, mngr)
}*/

/*func (s *S) TestManager(c *check.C) {
	mngr := nopManager{}
	Register("nope", mngr)
	config.Set("repo-manager", "nope")
	defer config.Unset("repo-manager")
	current := Manager()
	c.Assert(current, check.Equals, mngr)
}

func (s *S) TestManagerUnconfigured(c *check.C) {
	mngr := nopManager{}
	Register("nope", mngr)
	gitlab := nopManager{}
	Register("gandalf", gitlab)
	config.Unset("repo-manager")
	current := Manager()
	c.Assert(current, check.Equals, gandalf)
}

func (s *S) TestManagerUnknown(c *check.C) {
	config.Set("repo-manager", "something")
	defer config.Unset("repo-manager")
	current := Manager()
	c.Assert(current, check.FitsTypeOf, nopManager{})
}*/
