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
 
package gulpd

import (
//	"encoding/json"
//	"errors"
//	"io"
//	"log"
//	"time"
	"github.com/megamsys/gulp/activities"
	"github.com/megamsys/gulp/app"
)

type Handler struct {
	Provider string
	d        *Config
}

// NewHandler returns a new instance of handler.
func NewHandler(conf *Config) *Handler {
	h := &Handler{
		d: conf,
	}	

	return h
}

func (h *Handler) ServeAMQP(r *app.Requests) error {
	assembly, err := app.GetAssemblyWithComponents(r.Id)
	if err != nil {
		return err
	}

	di := app.ActionData{
		Assembly: assembly,
		Request: r,
			}
	
	p, err := activities.GetActivity(r.Category)	
	if err = p.Action(&di); err != nil {
			return err
		}	
	
	return nil
}
