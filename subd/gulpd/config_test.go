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
		provider = "gru"
		gru   = "github"
                gructl_tar = "https://s3-ap-southeast-1.amazonaws.com/megampub/gru-site/gructl.tar.gz"
		gru_tar = "https://github.com/megamsys/gru.git"
`, &cm); err != nil {
		c.Fatal(err)
	}

	c.Assert(cm.Provider, check.Equals, "gru")
	c.Assert(cm.GruGit, check.Equals, "github")
        c.Assert(cm.GructlTar, check.Equals, "https://s3-ap-southeast-1.amazonaws.com/megampub/gru-site/gructl.tar.gz")
	c.Assert(cm.GruTarball, check.Equals, "https://github.com/megamsys/gru.git")
}
