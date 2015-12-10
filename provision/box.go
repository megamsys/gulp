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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/loggers"
	_ "github.com/megamsys/gulp/loggers/file"
	_ "github.com/megamsys/gulp/loggers/queue"
	"github.com/megamsys/gulp/operations"
	"github.com/megamsys/gulp/repository"
	"gopkg.in/yaml.v2"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (

	// BoxSome indicates that there is atleast one box to deploy or delete.
	BoxSome BoxLevel = iota

	// BoxNone indicates that there are no boxes to deploy or delete but its parent can be.
	BoxNone
)

// Boxlevel represents the deployment level.
type BoxLevel int

var cnameRegexp = regexp.MustCompile(`^(\*\.)?[a-zA-Z0-9][\w-.]+$`)

// Box represents a provision unit. Can be a machine, container or anything
// IP-addressable.
type Box struct {
	Id         string
	CartonsId  string
	CartonId   string
	Level      BoxLevel
	Name       string
	DomainName string
	Tosca      string
	Repo       *repository.Repo
	Operations []*operations.Operate
	Status     Status
	Provider   string
	Commit     string
	Address    *url.URL
	Ip         string
	Cookbook   string
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
	return b.Status == StatusBootstrapping ||
		b.Status == StatusRunning ||
		b.Status == StatusBootstrapped ||
		b.Status == StatusStateup ||
		b.Status == StatusError ||
		b.Status == StatusStarted ||
		b.Status == StatusStopped ||
		b.Status == StatusRestarted
}

// Log adds a log message to the app. Specifying a good source is good so the
// user can filter where the message come from.
func (box *Box) Log(message, source, unit string) error {

	messages := strings.Split(message, "\n")
	logs := make([]loggers.Boxlog, 0, len(messages))
	for _, msg := range messages {
		if msg != "" {
			bl := loggers.Boxlog{
				Date:    time.Now().In(time.UTC),
				Message: msg,
				Name:    box.Name,
				Unit:    box.Id,
			}
			fmt.Println(bl.Message)
			logs = append(logs, bl)
		}
	}
	if len(logs) > 0 {
		a, err := loggers.Get(source)

		if err != nil {
			log.Errorf("fatal error, couldn't located the Logger %s", source)
			return err
		}

		Logger = a

		if initializableLogger, ok := Logger.(loggers.InitializableLogger); ok {
			err = initializableLogger.Notify(box.Name+"."+box.DomainName, logs)
			if err != nil {
				log.Errorf("fatal error, couldn't initialize the Logger %s", source)
				return err
			}
		}
		//_ = notify(box.Name+"."+box.DomainName, logs)
	}
	return nil
}
