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
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/repository"
	"github.com/megamsys/libgo/action"
)

func init() {
	repository.Register("github", githubManager{})
}

type githubManager struct{}

/**
* clone repository from github.com using url
**/
func (m githubManager) Clone(url string) error {

	actions := []*action.Action{
		&remove_old_file,
		&clone,
	}
	pipeline := action.NewPipeline(actions...)

	s := strings.Split(url, "/")[4]

	args := runActionsArgs{
		//	Writer:        w,
		url:      url,
		dir:      meta.MC.Dir,
		filename: strings.Split(s, ".")[0],
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute status pipeline for github %s - %s", url, err)
		return err
	}
	return nil

}
func (m githubManager) Initialize(url, tar_url string) error {

	actions := []*action.Action{
		&clone_tar,
		&make_dir,
		&un_tar,
		&remove_tar_file,
	}
	pipeline := action.NewPipeline(actions...)

	s := strings.Split(url, "/")[4]
	s1 := strings.Split(tar_url, "/")[6]
	args := runActionsArgs{
		//	Writer:        w,
		url:         url,
		tar_url:     tar_url,
		dir:         meta.MC.Dir,
		filename:    strings.Split(s, ".")[0],
		tarfilename: s1,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute status pipeline for github %s - %s", tar_url, err)
		return err
	}
	return nil

}
