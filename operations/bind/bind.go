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

package bind

import (
	//"fmt"
	//log "github.com/Sirupsen/logrus"
	//"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/operations"
)

func init() {
	operations.Register("bind", bindManager{})
}

type bindManager struct{}

func (m bindManager) Initialize(url string) error {
	return nil
}

/**
* clone repository from github.com using url
**/
/*func (m ciManager) Initialize(url string) error {

	actions := []*action.Action{
		&clone,
	}
	pipeline := action.NewPipeline(actions...)

	args := runActionsArgs{
	//	Writer:        w,
		Url:   url,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute status pipeline for github %s - %s", url, err)
		return err
	}
	return nil

}*/
