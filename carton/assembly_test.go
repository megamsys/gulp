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
 import (
	"testing"
  "fmt"

	"github.com/megamsys/gulp/carton/bind"
  "github.com/megamsys/gulp/db"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{
  A  *Ambly2
}

var _ = check.Suite(&S{})

type Ambly2 struct {
	Id            string         `json:"id"`
	Org_Id        string				 `json:"org_id"`
	Name          string         `json:"name"`
	Json_Claz     string         `json:"json_claz"`
	Tosca_Type    string         `json:"tosca_type"`
	Inputs        bind.JsonPairs `json:"inputs"`
	Outputs       bind.JsonPairs `json:"outputs"`
	Policies      []*Policy      `json:"policies"`
	Status        string         `json:"status"`
	Created_At    string         `json:"created_at"`
	Components    []string       `json:"components"`
}

func getBig2(id string,a *Ambly2) (*Ambly2, error) {
	t := db.TableInfo{
		Name: ASSEMBLYBUCKET,
		Pks: []string{"org_id","id"},
		Ccms: []string{},
		Query: map[string]string{"id": id},
	}
	if err := db.ReadWhere(t, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *S) TestFetchAssembly(c *check.C) {
  a := &Ambly2{}
  a,_ = getBig2("ASM4880170097874038597",a)
 fmt.Println(a)
	c.Assert(nil, check.IsNil)
}
