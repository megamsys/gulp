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
 
package run

import (
	"errors"
//	"fmt"
//	"os"
//	"os/user"
//	"path/filepath"
//	"reflect"
//	"strconv"
//	"strings"
//	"time"

	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/services/gulpd"
	"github.com/megamsys/gulp/services/httpd"
)

type Config struct {
	Meta    *meta.Config 	`toml:"meta"`
	Gulpd   *gulpd.Config   `toml:"gulpd"`
	HTTPD   *httpd.Config   `toml:"http"`

}

func (c Config) String() string {
	return (c.Meta.String() +
		c.Gulpd.String() + "\n" +
		c.HTTPD.String())
}

// NewConfig returns an instance of Config with reasonable defaults.
func NewConfig() *Config {
	c := &Config{}
	c.Meta = meta.NewConfig()

	c.Gulpd = gulpd.NewConfig()
	c.HTTPD = httpd.NewConfig()
	return c
}

// NewDemoConfig returns the config that runs when no config is specified.
func NewDemoConfig() (*Config, error) {
	c := NewConfig()
	return c, nil
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if c.Meta.Dir == "" {
		return errors.New("Meta.Dir must be specified")
	} else if c.Gulpd.AssemblyID == "" {
		return errors.New("Gulpd.AssemblyID must be specified")
	}
	return nil
}