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

// Package chefsolo implements a provisioner using Chef Solo.
package chefsolo

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	log "github.com/Sirupsen/logrus"
)

const (
	// DefaultFormat is the default output format of Chef.
	DefaultFormat = "doc"

	// DefaultLogLevel is the default log level of Chef.
	DefaultLogLevel = "info"
)

// Provisioner is a provisioner based on Chef Solo.
type Provisioner struct {
	RunList     []string
	Attributes  string
	Format      string
	LogLevel    string
	SandboxPath string
	RootPath    string
	Sudo        bool
}

func (p Provisioner) prepareJSON() error {
	log.Info("Preparing JSON data")
	data := "{}\n"
	if p.Attributes != "" {
		data = p.Attributes
	}
	return ioutil.WriteFile(path.Join(p.SandboxPath, "solo.json"), []byte(data), 0644)
}

func (p Provisioner) prepareSoloConfig() error {
	log.Info("Preparing Chef Solo config")
	data := fmt.Sprintf("cookbook_path \"%s\"\n", path.Join(p.RootPath, "cookbooks"))
	data += "ssl_verify_mode :verify_peer\n"
	return ioutil.WriteFile(path.Join(p.SandboxPath, "solo.rb"), []byte(data), 0644)
}

// PrepareFiles prepares Chef configuration data.
func (p Provisioner) PrepareFiles() error {
	if err := p.prepareJSON(); err != nil {
		return err
	}
	return p.prepareSoloConfig()
}

// Command returns the command string which will invoke the provisioner on the
// prepared machine.
func (p Provisioner) Command() []string {
	format := p.Format
	if format == "" {
		format = DefaultFormat
	}

	logLevel := p.LogLevel
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

	cmd := []string{
		"chef-solo",
		"--config", path.Join(p.RootPath, "solo.rb"),
		"--json-attributes", path.Join(p.RootPath, "solo.json"),
		"--format", format,
		"--log_level", logLevel,
	}

	if len(p.RunList) > 0 {
		cmd = append(cmd, "--override-runlist", strings.Join(p.RunList, ","))
	}

	if !p.Sudo {
		return cmd
	}
	return append([]string{"sudo"}, cmd...)
}
