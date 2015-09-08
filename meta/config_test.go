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
 
package meta

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

type S struct{}

var _ = check.Suite(&S{})

// Ensure the configuration can be parsed.
func (s *S) TestConfig_Parse(c *check.C) {
	// Parse configuration.
	var conf Config
	if _, err := toml.Decode(`
		[meta]
			debug = true
			hostname = "localhost"
			bind_address = ":9999"
			dir = "/var/lib/megam/gulp/meta"
			riak = "192.168.1.100:8087"
			api  = "https://api.megam.io/v2"
			amqp = "amqp://guest:guest@192.168.1.100:5672/"

`, &c); err != nil {
		//t.Fatal(err)
	}

	c.Assert(conf.Hostname, check.Equals, "localhost")
	c.Assert(conf.Riak, check.Equals, "192.168.1.100:8087")
	c.Assert(conf.Api, check.Equals, "https://api.megam.io/v2")

}
