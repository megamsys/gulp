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
	"text/tabwriter"
	"strconv"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/gulp/provision/chefsolo"
)

const (
	// DefaultAssemblyID.
	DefaultAssemblyID = "ASM00"

	// DefaultProvider is the default provisioner used by our engine.
	DefaultProvider = "chefsolo"

	// DefaultCookbook is the default cookbook used by chefsolo.
	DefaultCookbook = "megam_run"

	//DefaultRepository is the default repository for megam
	DefaultRepository = "github"

	//DefaultRepositoryPath is the default repository path by megam
	DefaultRepositoryPath = "https://github.com/megamsys/chef-repo.git"

  //default git release of chef-repo
  DefaultRepositoryTarPath = "https://github.com/megamsys/chef-repo/archive/0.9.tar.gz"

	DefaultHomeDir = "/var/lib/megam"
)

type Config struct {
	Enabled     bool   `toml:"enabled"`
	Name		       string		 `toml:"name"`
	CatsID			   string		 `toml:"cats_id"`
	CatID              string 		 `toml:"cat_id"`
	Provider 		   string 	     `toml:"provider"`
	Cookbook 		   string 	     `toml:"cookbook"`
	Repository	       string        `toml:"repository"`
	RepositoryPath     string        `toml:"repository_path"`
  RepositoryTarPath  string       `toml:"repository_tar_path"`
	HomeDir            string        `toml:"dir"`
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Gulpd", "green", "", "") + "\n"))
	b.Write([]byte("Enabled" + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("Name" + "\t" + c.Name + "\n"))
	b.Write([]byte("CatID" + "\t" + c.CatID + "\n"))
	b.Write([]byte("Provider" + "\t" + c.Provider + "\n"))
	b.Write([]byte("Cookbook" + "\t" + c.Cookbook + "\n"))
	b.Write([]byte("Repository" + "\t" + c.Repository + "\n"))
	b.Write([]byte("RepositoryPath" + "\t" + c.RepositoryPath +"\n"))
  b.Write([]byte("RepositoryTarPath" + "\t" + c.RepositoryTarPath ))
	b.Write([]byte("HomeDir" + "\t" + c.HomeDir ))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}

func NewConfig() *Config {
	return &Config{
		Enabled:     true,
		Name:	  			"",
		Provider: 			DefaultProvider,
		CatID:    			DefaultAssemblyID,
		Cookbook:    		DefaultCookbook,
		Repository:			DefaultRepository,
		RepositoryPath:     DefaultRepositoryPath,
    RepositoryTarPath:     DefaultRepositoryTarPath,
		HomeDir:						DefaultHomeDir,
	}
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[chefsolo.Repository] = c.Repository
	m[chefsolo.RepositoryPath] = c.RepositoryPath
  m[chefsolo.RepositoryTarPath] = c.RepositoryTarPath
	m[chefsolo.HomeDir] = c.HomeDir
	return m
}
