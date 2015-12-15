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
package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
)

type ReqOperator struct {
	Id string
}

// NewReqOperator returns a new instance of ReqOperator
func NewReqOperator(id string) *ReqOperator {
	return &ReqOperator{Id: id}
}

func (p *ReqOperator) Accept(r *MegdProcessor, cookbook string) error {
	c, err := p.Get(p.Id, cookbook)

	if err != nil {
		return err
	}

	md := *r

	log.Debugf(cmd.Colorfy(md.String(), "cyan", "", "bold"))

	return md.Process(c)
}

func (p *ReqOperator) Get(cat_id string, cookbook string) (*Carton, error) {
	a, err := Get(cat_id)
	if err != nil {
		return nil, err
	}

	ca, err := a.MkCarton(cookbook)
	if err != nil {
		return nil, err
	} else {
		ca.toBox(cookbook) //on success, make a carton2box if BoxLevel is BoxZero
	}

	if err != nil {
		return nil, err
	}
	return ca, nil
}

// MegdProcessor represents a single operation in Gulpd.
type MegdProcessor interface {
	Process(c *Carton) error
	String() string
}
