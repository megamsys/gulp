package handler

import (
	"encoding/json"
	"io/ioutil"
	libhttp "net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/megamsys/gulp/activities/docker"
	"github.com/megamsys/gulp/app"
)

const (
	LOGS       = "logs"
	NETWORKS   = "networks"
	MONITORING = "monitoring"
)

func Logs(w libhttp.ResponseWriter, req *libhttp.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
	}

	d := app.LogsInfo{}
	err = json.Unmarshal(body, &d)
	var dat = app.ActionData{
		DockerLogs: &d,
		DockerType: LOGS,
	}

	doc := docker.DockerActivity{}
	if err = doc.Action(&dat); err != nil {
		log.Error(err)
	}
}

func Networks(w libhttp.ResponseWriter, req *libhttp.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
	}

	d := app.NetworksInfo{}
	err = json.Unmarshal(body, &d)
	var dat = app.ActionData{
		DockerNetworks: &d,
		DockerType:     NETWORKS,
	}

	doc := docker.DockerActivity{}

	if err = doc.Action(&dat); err != nil {
		log.Error(err)
	}
}
