package http

/*import (
	"encoding/json"
	"fmt"
	//"io"
	"io/ioutil"
	libhttp "net/http"

	"github.com/megamsys/gulp/coordinator"
	"github.com/megamsys/gulp/global"
)*/

/*
 * megam docker handlers for endpoint - /docker/logs
 */
/*func DockerLogs(w libhttp.ResponseWriter, req *libhttp.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("error")
	}
	var jsonLogsData global.DockerLogsInfo
	err = json.Unmarshal(body, &jsonLogsData)
	coordinator.DockerLogs(jsonLogsData.ContainerId, jsonLogsData.ContainerName)
}*/

/*
 * megam docker handlers for endpoint - /docker/networks
 */
/*func DockerNetworks(w libhttp.ResponseWriter, req *libhttp.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("error")
	}
	var jsonNetworksData global.DockerNetworksInfo
	err = json.Unmarshal(body, &jsonNetworksData)
	coordinator.DockerNetworks(jsonNetworksData.Bridge, jsonNetworksData.ContainerId, jsonNetworksData.IpAddr, jsonNetworksData.Gateway)
}*/
