package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	libhttp "net/http"
	"github.com/megamsys/gulp/coordinator"
)

/*
 * megam docker endpoints used for docker logs and networks
 */

type Dockerlogs struct {
	ContainerId   string `json:"container_id"`
	ContainerName string `json:"container_name"`
}

type Dockernetworks struct {
}

func DockerLogs(w libhttp.ResponseWriter, req *libhttp.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(string(body))

	var jsonData Dockerlogs
	err = json.Unmarshal(body, &jsonData)
	coordinator.DockerLogs(jsonData.ContainerId, jsonData.ContainerName)

}

func DockerNetworks(w libhttp.ResponseWriter, req *libhttp.Request) {
	fmt.Fprint(w, "Hello, Network\n")
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println("Docker networks")
	fmt.Println(body)
}
