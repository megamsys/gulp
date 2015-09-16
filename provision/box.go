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
package provision

import (
	"github.com/megamsys/gulp/repository"
	"gopkg.in/yaml.v2"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var cnameRegexp = regexp.MustCompile(`^(\*\.)?[a-zA-Z0-9][\w-.]+$`)

// Boxlog represents a log entry.
type Boxlog struct {
	Date    time.Time
	Message string
	Source  string
	Name    string
	Unit    string
}

// Box represents a provision unit. Can be a machine, container or anything
// IP-addressable.
type Box struct {
	ComponentId string
	AssemblyId  string
	Name        string
	DomainName  string
	Tosca       string
	Commit      string
	Image       string
	Repo        repository.Repository
	Status      Status
	Provider    string
	Address     *url.URL
	Ip          string
}

func (b *Box) String() string {
	if d, err := yaml.Marshal(b); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

// GetName returns the name of the box.
func (b *Box) GetFullName() string {
	return b.Name + b.DomainName
}

// GetTosca returns the tosca type of the box.
func (b *Box) GetTosca() string {
	return b.Tosca
}

// GetIp returns the Unit.IP.
func (b *Box) GetIp() string {
	return b.Ip
}

// Available returns true if the unit is available. It will return true
// whenever the unit itself is available, even when the application process is
// not.
func (b *Box) Available() bool {
	return b.Status == StatusDeploying ||
		b.Status == StatusCreating ||
		b.Status == StatusError
}

// Log adds a log message to the app. Specifying a good source is good so the
// user can filter where the message come from.
func (box *Box) Log(message, source, unit string) error {
	messages := strings.Split(message, "\n")
	logs := make([]interface{}, 0, len(messages))
	for _, msg := range messages {
		if msg != "" {
			bl := Boxlog{
				Date:    time.Now().In(time.UTC),
				Message: msg,
				Source:  source,
				Name:    box.Name,
				Unit:    box.ComponentId,
			}
			logs = append(logs, bl)
		}
	}
	if len(logs) > 0 {
		//notify(bl.Name, logs)
	}
	return nil
}
