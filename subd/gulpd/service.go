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
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	nsqc "github.com/crackcomm/nsqueue/consumer"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	_ "github.com/megamsys/gulp/provision/chefsolo"
	"github.com/megamsys/libgo/events"
	constants "github.com/megamsys/libgo/utils"
	nsq "github.com/nsqio/go-nsq"
)

const (
	maxInFlight = 150
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	wg       sync.WaitGroup
	err      chan error
	Handler  *Handler
	Consumer *nsqc.Consumer
	Meta     *meta.Config
	Gulpd    *Config
}

type Sales struct {
	Id         string
	CustomerId string
	SellerId   string
	Price      string
	Created    string
}

// NewService returns a new instance of Service.
func NewService(c *meta.Config, d *Config) *Service {
	s := &Service{
		err:   make(chan error),
		Meta:  c,
		Gulpd: d,
	}
	s.Handler = NewHandler(s.Gulpd)
	c.MkGlobal()
	d.MkGlobal()
	carton.NewArgs("")
	return s
}

// Open starts the service
func (s *Service) Open() error {
	if err := s.setEventsWrap(); err != nil {
		return err
	}
	go func() error {
		log.Info("starting deployd agent service")
		if err := nsqc.Register(s.Meta.Name, "agent", maxInFlight, s.processNSQ); err != nil {
			return err
		}
		if err := nsqc.Connect(s.Meta.NSQd...); err != nil {
			return err
		}
		s.Consumer = nsqc.DefaultConsumer

		nsqc.Start(true)
		return nil
	}()
	if err := s.setProvisioner(); err != nil {
		return err
	}
	s.boot()
	return nil
}

func (s *Service) setEventsWrap() error {
	mi := make(map[string]map[string]string)
	mi[constants.META] = meta.MC.ToMap()
	return events.NewWrap(mi)
}

func (s *Service) processNSQ(msg *nsqc.Message) {
	p, err := carton.NewPayload(msg.Body)
	if err != nil {
		return
	}

	re, err := p.Convert()
	if err != nil {
		return
	}
	go s.Handler.serveNSQ(re)

	return
}

func (s *Service) boot() {
	go func() {
		b, err := (&carton.Payload{}).AsBytes("", s.Meta.CartonId,
			carton.BOOT, carton.STATE, time.Now().Local().Format(time.RFC822))
		if err != nil {
			return
		}
		s.processNSQ(&nsqc.Message{&nsq.Message{Body: b}})
	}()
}

// Close closes the underlying subscribe channel.
func (s *Service) Close() error {
	if s.Consumer != nil {
		s.Consumer.Stop()
	}

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
