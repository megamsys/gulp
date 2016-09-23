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

package run

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
	"os"
)

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	var cm Config
	u, _ := os.Getwd()
	if _, err := toml.DecodeFile(u+"/gulpd.conf", &cm); err != nil {
		fmt.Println(err.Error())
	}

	c.Assert(cm, check.NotNil)
	c.Assert(cm.Meta.Riak, check.DeepEquals, []string{"localhost:8087"})
	c.Assert(cm.Gulpd.Provider, check.Equals, "chefsolo")
}
