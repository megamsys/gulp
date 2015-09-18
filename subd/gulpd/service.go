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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"

	"sync"
	"time"
)

const leaderWaitTimeout = 30 * time.Second

const QUEUE = "cloudstandup"

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg      sync.WaitGroup
	err     chan error
	Handler *Handler

	Meta    *meta.Config
	Gulpd *Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, d *Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Gulpd: d,
	}
	s.Handler = NewHandler(s.Gulpd)
	c.MC() //an accessor.
	return s
}

// Open starts the service
func (s *Service) Open() error {

	log.Debug("starting gulpd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, QUEUE)
	if err != nil {
		return err
	}

	if swt, err := p.Sub(); err != nil {
		return err
	} else {
		if err = s.setProvisioner(); err != nil {
			return err
		}
		
		//before publish the queue, we need to verify assembly status
  		s.publishAMQP(QUEUE, &carton.Requests{Name: s.Gulpd.Name, CatId: s.Gulpd.CatID, CatType: carton.STATE, Action: provision.ReqStarted.string}.ToJson())

		go s.processQueue(swt)
	}

	return nil	
}

// processQueue continually drains the given queue  and processes the queue request
// to the appropriate handlers..
func (s *Service) processQueue(drain chan []byte) error {
	//defer s.wg.Done()
	for raw := range drain {
		p, err := carton.NewPayload(raw)
		if err != nil {
			return err
		}

		pc, err := p.Convert()
		if err != nil {
			return err
		}
		go s.Handler.serveAMQP(pc)
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

func (s *Service) publishAMQP(key string, json string) {
	factor, aerr := amqp.Factory()
	if aerr != nil {
		log.Errorf("Failed to get the queue instance: %s", aerr)
	}
	//s := strings.Split(key, "/")
	//pubsub, perr := factor.Get(s[len(s)-1])
	pubsub, perr := factor.Get(key)
	if perr != nil {
		log.Errorf("Failed to get the queue instance: %s", perr)
	}

	serr := pubsub.Pub([]byte(json))
	if serr != nil {
		log.Errorf("Failed to publish the queue instance: %s", serr)
	}
}

//this is an array, a property provider helps to load the provider specific stuff
func (s *Service) setProvisioner() {
	a, err := provision.Get(s.Meta.Provider)

	if err != nil {
		fmt.Errorf("fatal error, couldn't located the provisioner %s", s.Meta.Provider)
	}
	carton.Provisioner = a

	log.Debugf("Using %q provisioner. %q", s.Meta.Provider, a)
	if initializableProvisioner, ok := carton.Provisioner.(provision.InitializableProvisioner); ok {
		log.Debugf("Before initialization.")
		err = initializableProvisioner.Initialize(s.Meta.toMap())
		if err != nil {
			log.Errorf("fatal error, couldn't initialize the provisioner %s", s.Meta.Provider)
		} else {
			log.Debugf("%s Initialized", s.Meta.Provider)
		}
	}
	log.Debugf("After initialization.")

	if messageProvisioner, ok := carton.Provisioner.(provision.MessageProvisioner); ok {
		startupMessage, err := messageProvisioner.StartupMessage()
		if err == nil && startupMessage != "" {
			log.Debugf(startupMessage)
		} else {
			log.Debugf("------> " + err.Error())
		}
	}
}


