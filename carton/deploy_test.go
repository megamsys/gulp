/*
** copyright [2013-2016] [Megam Systems]
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
package carton

/*import (
	"testing"
  "fmt"
 "github.com/megamsys/gulp/carton/bind"
 "github.com/megamsys/gulp/provision"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}


type S struct{}

var _ = check.Suite(&S{})


/*func (s *S) TestBindService(c *check.C) {
 var z = make(bind.EnvVars,1)
 var envs = bind.EnvVar{
    Name: "port",
    Value: "8080",
    Endpoint: "",
  }
  z[0] = envs
 var box = provision.Box{
   Envs: z,
 }
 var x = DeployOpts{
   B: &box,
 }
  err :=BindService(&x)
  expected := `gulpd failed to apply the lifecle to the app "myapp": failure in app`
	c.Assert(err, check.Equals, expected)
}   // */
