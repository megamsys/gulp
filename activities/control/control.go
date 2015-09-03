package control

import (
	"log"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/handlers"
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
		go app.RebootApp(data.Assembly)
		break
	case "start":
		go app.StartApp(data.Assembly)
		break
	case "stop":
		go app.StopApp(data.Assembly)
		break
	case "restart":
		go app.RestartApp(data.Assembly)
		break
	case "redeploy":
		go app.BuildApp(data.Assembly)
		break
	}

	return nil
}
