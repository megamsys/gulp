package docker

import (
	"github.com/megamsys/gulp/activities"
	"github.com/megamsys/gulp/app"
)

func Init() {
	activities.RegisterActivities("docker", &DockerActivity{})
}

type DockerActivity struct{}

func (c *DockerActivity) Action(data *app.ActionData) error {

	switch data.DockerType {
	case "logs":
		go app.StreamLogs(data.DockerLogs)
	case "networks":
		go app.ConfigureNetworks(data.DockerNetworks)
	}

	return nil
}
