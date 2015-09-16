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
 
package github

import (
//	"encoding/json"
//	log "github.com/Sirupsen/logrus"
	//git "github.com/google/go-github/github"
//	"strings"
	"github.com/megamsys/gulp/repository"
)

func init() {
	repository.Register("github", githubManager{})
}

const endpointConfig = "git:api-server"

type githubManager struct{}

/**
* clone repository from github.com using url
**/
func (m githubManager) Clone(url string) error {
	
	return nil

}

