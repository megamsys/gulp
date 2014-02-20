package app

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/indykish/gulp/action"
	"github.com/indykish/gulp/db"
	"github.com/indykish/gulp/exec"
	"github.com/indykish/gulp/scm"
	"log"
	//"launchpad.net/goamz/aws"
	//"strconv"
	"strings"
)

const (
	keyremote_repo = "remote_repo="
	keylocal_repo  = "local_repo="
	keyproject     = "project="
)

var ErrAppAlreadyExists = errors.New("there is already an app with this name.")

func CommandExecutor(app *App) (action.Result, error) {
	var e exec.OsExecutor
	var b bytes.Buffer
	commandWords := strings.Fields(app.AppReqs.LCApply)
	fmt.Println(commandWords, len(commandWords))

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b); err != nil {
			return nil, err
		}
	}

	log.Printf("%s", b)
	return &app, nil
}

// insertApp is an action that inserts an app in the database in Forward and
// removes it in the Backward.
//
// The first argument in the context must be an App or a pointer to an App.
var startApp = action.Action{
	Name: "startapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		//err = conn.Apps().Insert(app)
		//if err != nil && strings.HasPrefix(err.Error(), "E11000") {
		//	return nil, ErrAppAlreadyExists
		//}
		//return &app, err
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		conn, err := db.Conn()
		if err != nil {
			log.Printf("Could not connect to the database: %s", err)
			return
		}
		log.Printf("App name is %s", app.Name)
		defer conn.Close()
		//conn.Apps().Remove(bson.M{"name": app.Name})
	},
	MinParams: 1,
}

// insertApp is an action that inserts an app in the database in Forward and
// removes it in the Backward.
//
// The first argument in the context must be an App or a pointer to an App.
var stopApp = action.Action{
	Name: "stopapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		//err = conn.Apps().Insert(app)
		//if err != nil && strings.HasPrefix(err.Error(), "E11000") {
		//	return nil, ErrAppAlreadyExists
		//}
		//return &app, err
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		conn, err := db.Conn()
		if err != nil {
			log.Printf("Could not connect to the database: %s", err)
			return
		}
		log.Printf("App name is %s", app.Name)
		defer conn.Close()
		//conn.Apps().Remove(bson.M{"name": app.Name})
	},
	MinParams: 1,
}

// insertApp is an action that inserts an app in the database in Forward and
// removes it in the Backward.
//
// The first argument in the context must be an App or a pointer to an App.
var buildApp = action.Action{
	Name: "buildapp",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var app App
		switch ctx.Params[0].(type) {
		case App:
			app = ctx.Params[0].(App)
		case *App:
			app = *ctx.Params[0].(*App)
		default:
			return nil, errors.New("First parameter must be App or *App.")
		}

		project, err := scm.Project()
		if err != nil {
			log.Printf("Could not find the project name in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the project name in gulp.conf file")
		}

		builder, err := scm.Builder()
		if err != nil {
			log.Printf("Could not find the builder in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the builder in gulp.conf file")

		}

		local_repo, err := scm.GetPath()
		if err != nil {
			log.Printf("Could not find the local repo  in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the local repo in gulp.conf file")
		}

		remote_repo, err := scm.GetRemotePath()
		if err != nil {
			log.Printf("Could not find the remote repo in gulp.conf file: %s", err)
			return nil, errors.New("Could not find the remote repo in gulp.conf file")
		}

		build_parms := fmt.Sprintf("%s/%s %s %s %s", builder, app.AppReqs.LCApply, keyproject+project, keylocal_repo+local_repo, keyremote_repo+remote_repo)

		app.AppReqs.LCApply = build_parms
		//err = conn.Apps().Insert(app)
		//if err != nil && strings.HasPrefix(err.Error(), "E11000") {
		//	return nil, ErrAppAlreadyExists
		//}
		//return &app, err
		return CommandExecutor(&app)
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.FWResult.(*App)
		conn, err := db.Conn()
		if err != nil {
			log.Printf("Could not connect to the database: %s", err)
			return
		}
		log.Printf("App name is %s", app.Name)
		defer conn.Close()
		//conn.Apps().Remove(bson.M{"name": app.Name})
	},
	MinParams: 1,
}

/*
// exportEnvironmentsAction exports megam's default environment variables in a
// new app. It requires a pointer to an App instance as the first parameter,
// and the previous result to be a *s3Env (it should be used after
var exportEnvironmentsAction = action.Action{
	Name: "export-environments",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		app := ctx.Params[0].(*App)
		err := app.Get()
		if err != nil {
			return nil, err
		}
		t, err := auth.CreateApplicationToken(app.Name)
		if err != nil {
			return nil, err
		}
		host, _ := config.GetString("host")
		envVars := []bind.EnvVar{
			{Name: "MEGAM_APPNAME", Value: app.Name},
			{Name: "MEGAM_HOST", Value: host},
			{Name: "MEGAM_API_KEY", Value: t.Token},
		}
		env, ok := ctx.Previous.(*s3Env)
		if ok {
			variables := map[string]string{
				"ENDPOINT":           env.endpoint,
				"LOCATIONCONSTRAINT": strconv.FormatBool(env.locationConstraint),
				"ACCESS_KEY_ID":      env.AccessKey,
				"SECRET_KEY":         env.SecretKey,
				"BUCKET":             env.bucket,
			}
			for name, value := range variables {
				envVars = append(envVars, bind.EnvVar{
					Name:         fmt.Sprintf("MEGAM_S3_%s", name),
					Value:        value,
					InstanceName: s3InstanceName,
				})
			}
		}
		err = app.setEnvsToApp(envVars, false, true)
		if err != nil {
			return nil, err
		}
		return ctx.Previous, nil
	},
	Backward: func(ctx action.BWContext) {
		app := ctx.Params[0].(*App)
		auth.DeleteToken(app.Env["MEGAM_API_KEY"].Value)
		if app.Get() == nil {
			s3Env := app.InstanceEnv(s3InstanceName)
			vars := make([]string, len(s3Env)+3)
			i := 0
			for k := range s3Env {
				vars[i] = k
				i++
			}
			vars[i] = "MEGAM_HOST"
			vars[i+1] = "MEGAM_APPNAME"
			vars[i+2] = "MEGAM_APIKEY"
			app.UnsetEnvs(vars, false)
		}
	},
	MinParams: 1,
}

*/
