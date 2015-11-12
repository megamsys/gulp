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

package repository

import (
	"fmt"
)

const (
	defaultManager = "github"
	CI             = "CI"
	CI_ENABLED     = "enabled"
	CI_TOKEN       = "token"
	CI_SOURCE      = "source"
	CI_USER        = "username"
	CI_URL         = "url"
	CI_TYPE        = "type"

	// IMAGE indicates that the repo is an image
	IMAGE = "image"

	// Git indicates that the repo is a GIT
	GIT = "git"
)

var managers map[string]InitializableRepository

/* Repository represents a repository managed by the manager. */
type Repo struct {
	Enabled  bool
	Type     string
	Token    string
	Source   string
	Url      string
	UserName string
	CartonId string
	BoxId    string
	OneClick string
}

func (r Repo) IsEnabled() bool {
	return r.Enabled
}

func (r Repo) GetType() string {
	return r.Type
}

func (r Repo) GetSource() string {
	return r.Source
}

func (r Repo) GetToken() string {
	return r.Token
}

func (r Repo) GetUrl() string {
	return r.Url
}

type Repository interface {
}

// RepositoryManager represents a manager of application repositories.
type InitializableRepository interface {
	Clone(url string) error
	Initialize(url, tar_url string) error
}

//type InitializableRepository interface { Initialize(url,tar_url string) error  }
// Get gets the named provisioner from the registry.
func Get(name string) (Repository, error) {
	p, ok := managers[name]
	if !ok {
		return nil, fmt.Errorf("unknown repository: %q", name)
	}
	return p, nil
}

// Manager returns the current configured manager, as defined in the
// configuration file.
func Manager(managerName string) InitializableRepository {
	if _, ok := managers[managerName]; !ok {
		managerName = "nop"
	}
	return managers[managerName]
}

// Register registers a new repository manager, that can be later configured
// and used.
func Register(name string, manager InitializableRepository) {
	if managers == nil {
		managers = make(map[string]InitializableRepository)
	}
	managers[name] = manager
}
