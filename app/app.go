package app

import (
	//"bytes"
	"encoding/json"
//	stderr "errors"
//	"fmt"
//	"github.com/globocom/config"
//	"github.com/indykish/gulp/action"
	"github.com/indykish/gulp/db"
//	"github.com/indykish/gulp/errors"
	//"io"
	//"os"
	"regexp"
	//"sort"
	//"strings"
	//"time"
)

var (
	nameRegexp  = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
	cnameRegexp = regexp.MustCompile(`^[a-zA-Z0-9][\w-.]+$`)
)

// App is the main type in megam. An app represents a real world application.
// This struct holds information about the app: its name, address, list of
// teams that have access to it, used platform, etc.
type App struct {
	//Env      map[string]bind.EnvVar
	Platform string `chef:"java"`
	Name     string
	Ip       string
	CName    string
	//	Units    []Unit
	Teams   []string
	Owner   string
	State   string
	Deploys uint

	//	hr hookRunner
}

// MarshalJSON marshals the app in json format. It returns a JSON object with
//the following keys: name, framework, teams, units, repository and ip.
func (app *App) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["name"] = app.Name
	result["platform"] = app.Platform
	result["teams"] = app.Teams
	//result["units"] = app.Units
	//result["repository"] = repository.ReadWriteURL(app.Name)
	result["ip"] = app.Ip
	result["cname"] = app.CName
	result["ready"] = app.State == "ready"
	return json.Marshal(&result)
}

// Get queries the database and fills the App object with data retrieved from
// the database. It uses the name of the app as filter in the query, so you can
// provide this field:
//
//     app := App{Name: "myapp"}
//     err := app.Get()
//     // do something with the app
func (app *App) Get() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	//return conn.Apps().Find(bson.M{"name": app.Name}).One(app)
	return nil
}

// StartsApp creates a new app.
//
// Starts the app :
//
//func StartApp(app *App, user *auth.User) error {
func StartApp(app *App) error {
	/*	teams, err := user.Teams()
		if err != nil {
			return err
		}
		if len(teams) == 0 {
			return NoTeamsError{}
		}
		if _, err := getPlatform(app.Platform); err != nil {
			return err
		}
		app.SetTeams(teams)
		app.Owner = user.Email
		if !app.isValid() {
			msg := "Invalid app name, your app should have at most 63 " +
				"characters, containing only lower case letters, numbers or dashes, " +
				"starting with a letter."
			return &errors.ValidationError{Message: msg}
		}
		actions := []*action.Action{&reserveUserApp, &createAppQuota, &insertApp}
		useS3, _ := config.GetBool("bucket-support")
		if useS3 {
			actions = append(actions, &createIAMUserAction,
				&createIAMAccessKeyAction,
				&createBucketAction, &createUserPolicyAction)
		}
		actions = append(actions, &exportEnvironmentsAction,
			&createRepository, &provisionApp)
		pipeline := action.NewPipeline(actions...)
		err = pipeline.Execute(app, user)
		if err != nil {
			return &AppCreationError{app: app.Name, Err: err}
		}
	*/
	return nil
}

/* setEnv sets the given environment variable in the app.
func (app *App) setEnv(env bind.EnvVar) {
	if app.Env == nil {
		app.Env = make(map[string]bind.EnvVar)
	}
	app.Env[env.Name] = env
	if env.Public {
		app.Log(fmt.Sprintf("setting env %s with value %s", env.Name, env.Value), "megam")
	}
}

// getEnv returns the environment variable if it's declared in the app. It will
// return an error if the variable is not defined in this app.
func (app *App) getEnv(name string) (bind.EnvVar, error) {
	var (
		env bind.EnvVar
		err error
		ok  bool
	)
	if env, ok = app.Env[name]; !ok {
		err = stderr.New("Environment variable not declared for this app.")
	}
	return env, err
}
*/
// GetName returns the name of the app.
func (app *App) GetName() string {
	return app.Name
}

// GetIp returns the ip of the app.
func (app *App) GetIp() string {
	return app.Ip
}

// GetPlatform returns the platform of the app.
func (app *App) GetPlatform() string {
	return app.Platform
}

func (app *App) GetDeploys() uint {
	return app.Deploys
}

/* Env returns app.Env
func (app *App) Envs() map[string]bind.EnvVar {
	return app.Env
}*/
