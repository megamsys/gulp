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
//	"bufio"
//	"bytes"
//	"io"
	log "github.com/golang/glog"
//	"os"
//	"strconv"
//	"strings"
	"sync"
	"time"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/meta"
)

const leaderWaitTimeout = 30 * time.Second

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg      sync.WaitGroup
	err     chan error
	Handler *Handler

	Meta    *meta.Config
	Gulpd   *Config
}

// NewService returns a new instance of Service.
func NewService(c meta.Config, d Config) (*Service, error) {
	s := &Service{
		err:     make(chan error),
		Meta:    &c,
		Gulpd:   &d,
	}
	s.Handler = NewHandler(s.Gulpd)
	return s, nil
}

// Open starts the service
func (s *Service) Open() error {
	log.Info("Starting gulpd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, s.Gulpd.AssemblyID)
	if err != nil {
		log.Error("Couldn't establish an amqp (%s): %s", s.Meta, err.Error())
	}

	ch, err := p.Sub()

	for raw := range ch {
		p, err := app.NewPayload(raw)
		if err != nil {
			return err
		}
		req, rerr := p.Convert() 
		if rerr != nil {
			return rerr
		}	
		go s.Handler.ServeAMQP(req)
	}

	return nil
}

// Close closes the underlying subscribe channel.
func (s *Service) Close() error {
	/*save the subscribe channel and close it.
	  don't know if the amqp has Close method ?
	  	if s.chn != nil {
	  		return s.chn.Close()
	  	}
	*/
	s.wg.Wait()
	return nil
}

// Err returns a channel for fatal errors that occur on the listener.
func (s *Service) Err() <-chan error { return s.err }
