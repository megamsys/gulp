package app

import (
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/policies"
	"github.com/megamsys/gulp/global"
)


// StartsApp creates a new app.
//
// Starts the app :
func RestartApp(app *policies.AssemblyResult) error {
	actions := []*action.Action{&restartApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}


func StartApp(app *policies.AssemblyResult) error {
	actions := []*action.Action{&startApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func StopApp(app *policies.AssemblyResult) error {
	actions := []*action.Action{&stopApp}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}


func RestartComponent(app *global.Component) error {
	actions := []*action.Action{&restartComponent}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}


func StartComponent(app *global.Component) error {
	actions := []*action.Action{&startComponent}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func StopComponent(app *global.Component) error {
	actions := []*action.Action{&stopComponent}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func RestartContainer(app *global.Component) error {
	actions := []*action.Action{&restartContainer}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}


func StartContainer(app *global.Component) error {
	actions := []*action.Action{&startContainer}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func StopContainer(app *global.Component) error {
	actions := []*action.Action{&stopContainer}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

func LogFile(app *global.Component) error {
	actions := []*action.Action{&logFile}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

