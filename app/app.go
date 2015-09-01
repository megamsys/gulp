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
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/state/provisioner/chefsolo"
	"github.com/megamsys/gulp/handlers"
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




func StartApp(app *handlers.AssemblyWithComponents) error {
	actions := []*action.Action{&startApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func StopApp(app *handlers.AssemblyWithComponents) error {
	actions := []*action.Action{&stopApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func RebootApp(app *handlers.AssemblyWithComponents) error {
	actions := []*action.Action{&rebootApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}



func RestartApp(app *handlers.AssemblyWithComponents) error {
	actions := []*action.Action{&restartApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func BuildApp(app *handlers.Component) error {
	actions := []*action.Action{&buildApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}
