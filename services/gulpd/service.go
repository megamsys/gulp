package gulpd

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

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
	Gulpd *gulpd.Config
}

// NewService returns a new instance of Service.
func NewService(c meta.Config, d gulpd.Config) (*Service, error) {
	if err != nil {
		return nil, err
	}

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
	log.Infof("Starting gulpd service")

	p, err := amqp.NewRabbitMQ(s.Meta.AMQP, s.Gulpd.AssemblyID)
	if err != nil {
		log.Errorf("Couldn't establish an amqp (%s): %s", s.Meta, err.Error())
	}

	ch, err := p.Sub()

	for raw := range ch {
		p, err := app.NewPayload(raw)
		if err != nil {
			return err
		}
		go s.Handler.serveAMQP(p.Convert())
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
