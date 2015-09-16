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
)

const (
	// DefaultAssemblyID.
	DefaultAssemblyID = "ASM00"
)

type Config struct {
	Enabled       string   			`toml:"enabled"`
	Name		  string			`toml:"name"`
	CatID         string 			`toml:"cat_id"`	
}

func (c Config) String() string {
	table := cmd.NewTable()
	table.AddRow(cmd.Row{cmd.Colorfy("Config:", "white", "", "bold"), cmd.Colorfy("Activity", "green", "", "")})
	table.AddRow(cmd.Row{"Enabled", c.Enabled})
	table.AddRow(cmd.Row{"Name", c.Name})	
	table.AddRow(cmd.Row{"CatID", c.CatID})	
	table.AddRow(cmd.Row{"", ""})
	return table.String() 
}

func NewConfig() *Config {
	return &Config{
		Enabled:     "true",
		Name:	  "",
		CatID:    DefaultAssemblyID,	
	}
}
