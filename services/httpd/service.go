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
	"fmt"
	log "github.com/golang/glog"
	"net"
	"net/http"
//	"os"
	"strings"
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
func NewService(c meta.Config, h Config) *Service {
	s := &Service{
		addr: h.BindAddress,
		err:  make(chan error),
		Handler: NewHandler(),
	}
	return s
}

// Open starts the service
func (s *Service) Open() error {
	log.Info("Starting HTTP service")

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	log.Info("Listening on HTTP:", listener.Addr().String())
	s.ln = listener

	// Begin listening for requests in a separate goroutine.
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

// Addr returns the listener's address. Returns nil if listener is closed.
func (s *Service) Addr() net.Addr {
	if s.ln != nil {
		return s.ln.Addr()
	}
	return nil
}

// serve serves the handler from the listener.
func (s *Service) serve() {
	// The listener was closed so exit
	// See https://github.com/golang/go/issues/4373
	err := http.Serve(s.ln, s.Handler)
	if err != nil && !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("listener failed: addr=%s, err=%s", s.Addr(), err)
	}
}
