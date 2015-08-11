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
package coordinator

import (
	"encoding/json"

	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/gulp/policies"
	"github.com/tsuru/config"
)

type Coordinator struct {
	//RequestHandler(f func(*Message), name ...string) (Handler, error)
	//EventsHandler(f func(*Message), name ...string) (Handler, error)
}

type Message struct {
	Id     string `json:"Id"`
	Action string `json:"Action"`
	Args   string `json:"Args"`
}

func Handler(chann []byte) {

	m := &Message{}
	parse_err := json.Unmarshal(chann, &m)
	log.Info(parse_err)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}
	log.Info("============Request entry===========")
	apprequest := global.AppRequest{Id: m.Id}
	req, err := apprequest.Get(m.Id)
	log.Info(req)
	if err != nil {
		log.Error("Error: Riak didn't cooperate->:\n%s.", err)
		return
	}

	assembly := global.Assembly{Id: req.AppId}
	asm, err := assembly.GetAssemblyWithComponents(req.AppId)
	if err != nil {
		log.Error("Error: Riak didn't cooperate-->:\n%s.", err)
		return
	}

	comp := global.Component{Id: asm.Components[0].Id}
	com, err := comp.Get(asm.Components[0].Id)
	if err != nil {
		log.Error("Error: Riak didn't cooperate--->:\n%s.", err)
		return
	}

	switch req.Action {
	case "reboot":
		log.Info("============Reboot entry==========")
		go app.RebootApp(asm)
		break
	case "start":
		log.Info("============Start entry===========")
		go app.StartApp(asm)
		break
	case "stop":
		log.Info("============Stop entry============")
		go app.StopApp(asm)
		break
	case "restart":
		log.Info("============Restart entry============")
		go app.RestartApp(asm)
		break
	case "redeploy":
		log.Info("=============Redeploying=============")
		go app.BuildApp(com)
		break
		/*case "componentstart":
		log.Info("============Component Start entry======")
			go app.StartComponent(com)
			break
		case "componentstop":
		log.Info("============Component Stop entry======")
			go app.StopComponent(com)
			break
		case "componentrestart":
		log.Info("============Component Restart entry======")
			go app.RestartComponent(com)
			break	  */
	}
}

/**
** It handles all events from megam engine
**/
func EventsHandler(chann []byte) {
	m := &Message{}
	parse_err := json.Unmarshal(chann, &m)
	log.Info(parse_err)
	if parse_err != nil {
		log.Error("Error: Message parsing error:\n%s.", parse_err)
		return
	}

	switch m.Action {
	case "build":
		log.Info("============Build entry============")
		comp := global.Component{Id: m.Id}
		com, err := comp.Get(m.Id)
		if err != nil {
			log.Error("Error: Riak didn't cooperate:\n%s.", err)
			return
		}
		go app.BuildApp(com)
		break
	}
}

func PolicyHandler() {
	log.Info("==>Policy Handler entry==")
	id, err := config.GetString("id")
	if err != nil {
		return
	}

	assembly := global.Assembly{Id: id}
	asm, asmerr := assembly.GetAssemblyWithComponents(id)
	if asmerr != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", asmerr)
		return
	}
	policies.ApplyPolicies(asm)
}

/*
* handlers to call the docker log stream and network setting action
*
 */

func DockerLogs(container_id string, container_name string) {
	log.Info("===>Docker Logs Entry<===")
	log := global.DockerLogsInfo{ContainerId: container_id, ContainerName: container_name}
	go app.StreamLogs(&log)
}

func DockerNetworks(bridge string, container_id string, ip_addr string, gateway string) {
	log.Info("===>Docker Networks Entry<===")
	network := global.DockerNetworksInfo{Bridge: bridge, ContainerId: container_id, IpAddr: ip_addr, Gateway: gateway}
	go app.ConfigureNetworks(&network)
}
