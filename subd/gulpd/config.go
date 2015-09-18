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
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/gulp/provision/chefsolo"
)

const (
	// DefaultAssemblyID.
	DefaultAssemblyID = "ASM00"	
	
	// DefaultProvider is the default provisioner used by our engine.
	DefaultProvider = "chefsolo"
	
	//DefaultRepository is the default repository for megam
	DefaultRepository = "github" 
	
	//DefaultRepositoryPath is the default repository path by megam
	DefaultRepositoryPath = "https://github.com/megamsys/chef-repo.git"
)

type Config struct {
	Name		       string		 `toml:"name"`
	CatID              string 		 `toml:"cat_id"`	
	Provider 		   string 	     `toml:"provider"`
	Repository	       string        `toml:"repository"`
	RepositoryPath     string        `toml:"repository_path"`
}

func (c Config) String() string {
	table := cmd.NewTable()
	table.AddRow(cmd.Row{cmd.Colorfy("Config:", "white", "", "bold"), cmd.Colorfy("Activity", "green", "", "")})
	table.AddRow(cmd.Row{"Name", c.Name})	
	table.AddRow(cmd.Row{"CatID", c.CatID})	
	table.AddRow(cmd.Row{"", ""})
	return table.String() 
}

func NewConfig() *Config {
	return &Config{
		Name:	  			"",
		Provider: 			DefaultProvider,
		CatID:    			DefaultAssemblyID,	
		Repository:			DefaultRepository,
		RepositoryPath:     DefaultRepositoryPath,
	}
}

//convert the config to just a map.
func (c Config) toMap() map[string]string {
	m := make(map[string]string)
	m[chefsolo.Repository] = c.Repository
	m[chefsolo.RepositoryPath] = c.RepositoryPath
	return m
}
