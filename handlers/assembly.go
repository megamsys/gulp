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
package handler

import (
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/db"
	//"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/log"
	"github.com/megamsys/megamd/provisioner"
)

type assembly struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	JsonClaz     string          `json:"json_claz"`
	ToscaType    string          `json:"tosca_type"`
	Requirements []*KeyValuePair `json:"requirements"`
	Policies     []*Policy       `json:"policies"`
	Inputs       []*KeyValuePair `json:"inputs"`
	Operations   []*Operations   `json:"operations"`
	Outputs      []*KeyValuePair `json:"outputs"`
	Status       string          `json:"status"`
	CreatedAt    string          `json:"created_at"`
}

/*type Assembly struct {
	assembly
	Components []string `json:"components"`
}*/

type DeepAssembly struct {
	assembly
	Components []*Component
}

func Deep(asmId string) (*DeepAssembly, error) {
	log.Debugf("[global] Get assembly_w_components %s", asmId)
	d := &DeepAssembly{}
	if conn, err := db.Conn("assembly"); err != nil {
		return d, err
	}

	if err := conn.FetchStruct(asmId, d); err != nil {
		return d, ferr
	}

	d.digMore()
	defer conn.Close()
	return result, nil
}

func (asm *DeepAssembly) digMore() error {
	for i := range asm.Components {
		if len(strings.TrimSpace(asm.Components[i])) > 1 {
			comp := NewComponent(asm.Components[i])
			if err := comp.Get(comp.Id); err != nil {
				log.Errorf("Failed to get component %s from riak: %s.", comp.Id, err.Error())
				return err
			}
			asm.ComponentsMap[comp.Id] = comp
		}
	}
}
