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
	"github.com/megamsys/libgo/action"
	_"github.com/megamsys/gulp/provision/chefsolo"
	"sync"
	"time"
	"net"
	"encoding/json"
)

const leaderWaitTimeout = 30 * time.Second

const QUEUE = "cloudstandup"

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg      sync.WaitGroup
	err     chan error
	Handler *Handler

	Meta    *meta.Config
	Gulpd 	*Config
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, d *Config) *Service {
	s := &Service{
		err:     make(chan error),
		Meta:    c,
		Gulpd:   d,
	}
	s.Handler = NewHandler(s.Gulpd)
	c.MC() //an accessor.
	return s
}

// Open starts the service
func (s *Service) Open() error {

	log.Debug("starting gulpd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, s.Gulpd.Name)
	if err != nil {
		return err
	}

	if swt, err := p.Sub(); err != nil {
		return err
	} else {
		if err = s.setProvisioner(); err != nil {
			return err
		}
		
		if err = s.updateStatusPipeline(); err != nil {
			return err
		}
		
		go s.processQueue(swt)
	}

	return nil	
}

// GetLocalIP returns the non loopback local IP of the host
func (s *Service) GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
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
		go s.Handler.serveAMQP(pc, s.Gulpd.Cookbook)
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


//this is an array, a property provider helps to load the provider specific stuff
func (s *Service) setProvisioner() error {
	var err error

	if carton.Provisioner, err = provision.Get(s.Gulpd.Provider); err != nil {
		return err
	}

	log.Debugf("configuring %s provisioner", s.Gulpd.Provider)
	if initializableProvisioner, ok := carton.Provisioner.(provision.InitializableProvisioner); ok {
		err = initializableProvisioner.Initialize(s.Gulpd.toMap())
		if err != nil {
			return fmt.Errorf("unable to initialize %s provisioner\n --> %s", s.Gulpd.Provider, err)
		} else {
			log.Debugf("%s initialized", s.Gulpd.Provider)
		}
	}

	if messageProvisioner, ok := carton.Provisioner.(provision.MessageProvisioner); ok {
		startupMessage, err := messageProvisioner.StartupMessage()
		if err == nil && startupMessage != "" {
			log.Infof(startupMessage)
		}
	}
	return nil
}

//1. &updateStatus in Riak - Bootstrapped..
//2. &publishStatus in publish the bootstrapped message to cloudstandup queue 
func (s *Service) updateStatusPipeline() error {
	actions := []*action.Action{		
		&updateStatusInRiak,
		&updateIPInRiak,
		&publishStatus,
	}
	pipeline := action.NewPipeline(actions...)	

	err := pipeline.Execute(s)
	if err != nil {
		log.Errorf("error on execute update pipeline for service %s - %s", "Gulpd", err)
		return err
	}
	return nil
}

//just publish a message bootstrapped to the service.
func (s *Service) pubStatus(status provision.Status) error {
	log.Debugf("publish bootstrapped service %s to %s", "Gulpd", status.String())

	p, err := amqp.NewRabbitMQ(meta.MC.AMQP, QUEUE)
	if err != nil {
		return err
	}
	
    //before publish the queue, we need to verify assembly status
    jsonMsg, err := json.Marshal(
		carton.Requests{
			CatId: 		s.Gulpd.CatsID, 
			Action:     status.String(),
			Category:   carton.STATE,
			CreatedAt:  time.Now().String(),
		})

	if err != nil {
		return err
	}	

	if err := p.Pub(jsonMsg); err != nil {
		return err
	}
	return nil

}


