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
package app

import (
	"github.com/megamsys/gulp/activities/state/provisioner/chefsolo"
	"github.com/megamsys/libgo/action"
)

/**
** state up the virtual machine
**/
func StateUP(app *chefsolo.Provisioner) error {
	actions := []*action.Action{&stateup}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: "", Err: err}
	}
	return nil
}

func StartApp(app *AssemblyWithComponents) error {
	actions := []*action.Action{&startApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func StopApp(app *AssemblyWithComponents) error {
	actions := []*action.Action{&stopApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func RebootApp(app *AssemblyWithComponents) error {
	actions := []*action.Action{&rebootApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func RestartApp(app *AssemblyWithComponents) error {
	actions := []*action.Action{&restartApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func BuildApp(app *Component) error {
	actions := []*action.Action{&buildApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

/*
* Docker logs stream which links docker logs to megam docker file for heka to read
*
 */

func StreamLogs(logs *LogsInfo) error {
	actions := []*action.Action{&streamLogs}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(logs)
	if err != nil {
		return &AppLifecycleError{app: logs.ContainerName, Err: err}
	}
	return nil
}

/*
* Docker networks configuration to setting up public ip
*
 */

func ConfigureNetworks(networks *NetworksInfo) error {

	actions := []*action.Action{&configureNetworks}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(networks)
	if err != nil {
		return &AppLifecycleError{app: networks.ContainerId, Err: err}
	}
	return nil
}
