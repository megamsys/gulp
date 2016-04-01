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
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/provision/chefsolo"
	"github.com/megamsys/libgo/cmd"
)

const (

	// DefaultProvider is the default provisioner used by our engine.
	DefaultProvider = provision.CHEFSOLO

	// DefaultCookbook is the default cookbook used by chefsolo.
	DefaultCookbook = "megam_run"

	//DefaultChefRepoGit is the default git for the chef-repo
	DefaultChefRepoGit = "https://github.com/megamsys/chef-repo.git"

	//DefaultChefTarball is the stable latest tar version
	DefaultChefTarball = "https://github.com/megamsys/chef-repo/archive/0.94.tar.gz"
)

var MC *Config

type Config struct {
	Enabled         bool   `toml:"enabled"`
	Provider        string `toml:"provider"`
	Cookbook        string `toml:"cookbook"`
	ChefRepoGit     string `toml:"chefrepo"`
	ChefRepoTarball string `toml:"chefrepo_tarball"`
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Gulpd", "green", "", "") + "\n"))
	b.Write([]byte("Enabled" + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("Provider" + "\t" + c.Provider + "\n"))
	b.Write([]byte("Cookbook" + "\t" + c.Cookbook + "\n"))
	b.Write([]byte("ChefRepoGit" + "\t" + c.ChefRepoGit + "\n"))
	b.Write([]byte("ChefRepoTarball" + "\t" + c.ChefRepoTarball + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}

func NewConfig() *Config {
	return &Config{
		Enabled:         true,
		Provider:        DefaultProvider,
		Cookbook:        DefaultCookbook,
		ChefRepoGit:     DefaultChefRepoGit,
		ChefRepoTarball: DefaultChefTarball,
	}
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[chefsolo.NAME] = meta.MC.Name
	m[chefsolo.CHEFREPO_GIT] = c.ChefRepoGit
	m[chefsolo.CHEFREPO_TARBALL] = c.ChefRepoTarball
	m[chefsolo.CHEFREPO_COOKBOOK] = c.Cookbook
	return m
}

func (c *Config) MkGlobal() {
	MC = c
}
