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
	"github.com/megamsys/gulp/operations"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/libgo/action"
)

func init() {
	operations.Register("bind", bindManager{})
}

const (
	APP             = "app"
	SERVICE         = "service"
	ASSEMBLIESINDEX = "assemblies_bin"
)

type bindManager struct{}


func (m bindManager) Initialize(url string) error {
	return nil
}

func (m bindManager) Apply(asm []*operations.Operate,envs []bind.EnvVar) (string, error) {
	for k := range asm {
		if asm[k].OperationType == "bind" {
		err := uploadENVVariables(asm, envs)
		 if err != nil {
			 log.Errorf("error on execute create pipeline for Set envs")
			 return "", err
		 }
		}
	}
	return "", nil
}

func uploadENVVariables(oprts []*operations.Operate,envs []bind.EnvVar) error {
	actions := []*action.Action{
		&setEnvs,
		&restartGulp,
	}
	pipeline := action.NewPipeline(actions...)
	args := runBindActionsArgs{
		envs:           envs,
		operations:        oprts,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute create pipeline for Set envs - %s", err)
		return err
	}
	return nil
}
