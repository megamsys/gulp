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
	"github.com/megamsys/gulp/provision/gru"
	"github.com/megamsys/libgo/cmd"
)

const (

	// DefaultProvider is the default provisioner used by our engine.
	DefaultProvider = provision.GRU

	// DefaultCookbook is the default cookbook used by gru.
	DefaultGructlTar = "https://s3-ap-southeast-1.amazonaws.com/megampub/gru-site/gructl.tar.gz"

	//DefaultChefRepoGit is the default git for the chef-repo
	DefaultGruGit = "https://github.com/megamsys/gru.git"

	//DefaultChefTarball is the stable latest tar version
	DefaultGruTarball = "https://github.com/megamsys/gru/archive/0.1.tar.gz"
)

var MC *Config

type Config struct {
	Enabled         bool   `toml:"enabled"`
	Provider        string `toml:"provider"`
	GructlTar       string `toml:"gructl_tar"`
	GruGit     string `toml:"gru"`
	GruTarball string `toml:"gru_tarball"`
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Gulpd", "green", "", "") + "\n"))
	b.Write([]byte("Enabled" + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("Provider" + "\t" + c.Provider + "\n"))
	b.Write([]byte("GructlTar" + "\t" + c.GructlTar + "\n"))
	b.Write([]byte("GruGit" + "\t" + c.GruGit + "\n"))
	b.Write([]byte("GruTarball" + "\t" + c.GruTarball + "\n"))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}

func NewConfig() *Config {
	return &Config{
		Enabled:         true,
		Provider:        DefaultProvider,
		GructlTar:        DefaultGructlTar,
		GruGit:     DefaultGruGit,
		GruTarball: DefaultGruTarball,
	}
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[gru.NAME] = meta.MC.Name
	m[gru.GRU_GIT] = c.GruGit
	m[gru.GRU_TARBALL] = c.GruTarball
	m[gru.GRUCTL_TAR] = c.GructlTar
	return m
}

func (c *Config) MkGlobal() {
	MC = c
}
