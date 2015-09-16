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
	"github.com/megamsys/gulp/carton"
)

type Handler struct {
	d            *Config
	EventChannel chan bool
}

func NewHandler(c *Config) *Handler {
	return &Handler{d: c}

}

func (h *Handler) serveAMQP(r *carton.Requests) error {
	if _, err := carton.ParseRequest(r.Name); err != nil {
		if rp := carton.NewReqOperator("nil"); rp != nil {
			return rp.Accept(nil)
		}
	}
	return nil
}

