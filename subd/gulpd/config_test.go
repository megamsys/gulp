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
	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

// Ensure the configuration can be parsed.
func (s *S) TestDeploydConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Config
	if _, err := toml.Decode(`
		provider = "chefsolo"
		repository   = "github"
		repository_path = "https://github.com/megamsys/chef-repo.git"
`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Provider, check.Equals, "chefsolo")
	c.Assert(cm.Repository, check.Equals, "github")
	c.Assert(cm.RepositoryPath, check.Equals, "https://github.com/megamsys/chef-repo.git")
}
