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

package control

import (
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/activities"
)

/**
**Control Activity register function
**This function register to activities container
**/
func Init() {
	activities.RegisterActivities("control", &ControlActivity{})
}

type ControlActivity struct{}

const (
	assemblyBucket = "assembly"
	comBucket      = "components"
)

func (c *ControlActivity) Action(data *app.ActionData) error {

	switch data.Request.Action {
	case "reboot":
	//	go app.RebootApp(data.Assembly)
		break
	case "start":
	//	go app.StartApp(data.Assembly)
		break
	case "stop":
	//	go app.StopApp(data.Assembly)
		break
	case "restart":
	//	go app.RestartApp(data.Assembly)
		break
	case "redeploy":
	//	go app.BuildApp(data.Assembly)
		break
	}

	return nil
}
