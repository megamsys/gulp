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
 
package httpd

import (
//	"encoding/json"
//	"io"
//	"net/http/httptest"
//	"reflect"
//	"regexp"
//	"time"
	"gopkg.in/check.v1"
)


// NewHandler represents a test wrapper for httpd.Handler.
//type Handler struct {
//	*httpd.Handler
//}

// NewHandler returns a new instance of Handler.
func NewHAndler() *Handler {
	h := &Handler{}
	h.Version = "0.0.0"
	return h
}

func (s *S) SetUpSuite(c *check.C) {
	h := NewHAndler()
	c.Assert(h.Version, check.Equals, "0.0.0")
}
