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

package httpd

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

type S struct{}

var _ = check.Suite(&S{})

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var h Config
	if _, err := toml.Decode(`
enabled = true
bind_address = ":8080"
`, &h); err != nil {
	}

	c.Assert(h.BindAddress, check.Equals, "locahost")
}
