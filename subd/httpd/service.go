/*
** Copyright [2013-2016] [Megam Systems]
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
	"fmt"
	"net"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/meta"
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	ln   net.Listener
	addr string
	err  chan error

	Handler *Handler
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, h *Config) (*Service, error) {
	s := &Service{
		addr:    h.BindAddress,
		err:     make(chan error),
		Handler: NewHandler(c, h),
	}
	return s, nil
}

// Open starts the service
func (s *Service) Open() error {
	log.Info("Starting HTTP service")

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	log.Info("Listening on HTTP:", s.addr)
	s.ln = listener
	go s.serve()
	return nil
}

// Close closes the underlying listener.
func (s *Service) Close() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error { return s.err }

// serve serves the handler from the listener.
func (s *Service) serve() {
	err := http.Serve(s.ln, s.Handler)
	if err != nil && !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("HTTP failed: %s, %s", s.addr, err)
	}
}
