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

package state

import (
	"gopkg.in/check.v1"
	"testing"
	"github.com/megamsys/gulp/app"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}


var _ = check.Suite(&S{})

func (s *S) TestStateUP(c *check.C) {
		di := app.ActionData{
			Assembly: &app.AssemblyWithComponents{},
			Request: &app.Requests{},
		}
		sat := StateActivity{}
		serr := sat.Action(&di)
		c.Assert(serr, check.IsNil)		
}