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

package run

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/subd/gulpd"
	"github.com/megamsys/gulp/subd/httpd"
)

// Server represents a container for the metadata and storage data and services.
// It is built using a Config and it manages the startup and shutdown of all
// services in the proper order.
type Server struct {
	version string // Build version

	err     chan error
	closing chan struct{}

	Services []Service

	// Profiling
	CPUProfile string
	MemProfile string
}

// NewServer returns a new instance of Server built from a config.
func NewServer(c *Config, version string) (*Server, error) {
	// Construct base meta store and data store.
	s := &Server{
		version: version,
		err:     make(chan error),
		closing: make(chan struct{}),
	}
	// Append services.
	s.appendGulpdService(c.Meta, c.Gulpd)
	s.appendHTTPDService(c.Meta, c.HTTPD)
	return s, nil
}

func (s *Server) appendGulpdService(c *meta.Config, d *gulpd.Config) {
	srv := gulpd.NewService(c, d)
	e := *d
	if !e.Enabled {
		log.Warn("skip gulpd service.")
		return
	}
	s.Services = append(s.Services, srv)
}

func (s *Server) appendHTTPDService(c *meta.Config, h *httpd.Config) {
	e := *h
	if !e.Enabled {
		log.Warn("skip httpd service.")
		return
	}

	srv, err := httpd.NewService(c, h)
	if err != nil {
		return
	}
	srv.Handler.Version = s.version
	s.Services = append(s.Services, srv)
}

// Err returns an error channel that multiplexes all out of band errors received from all services.
func (s *Server) Err() <-chan error { return s.err }

// Open opens the meta and data store and all services.
func (s *Server) Open() error {
	if err := func() error {
		startProfile(s.CPUProfile, s.MemProfile)
		for _, service := range s.Services {
			if err := service.Open(); err != nil {
				return fmt.Errorf("open service: %s", err)
			}
		}
		return nil
	}(); err != nil {
		s.Close()
		return err
	}
	return nil
}

// Close shuts down the meta and data stores and all services.
func (s *Server) Close() error {
	stopProfile()

	for _, service := range s.Services {
		service.Close()
	}

	close(s.closing)
	return nil
}

// monitorErrorChan reads an error channel and resends it through the server.
func (s *Server) monitorErrorChan(ch <-chan error) {
	for {
		select {
		case err, ok := <-ch:
			if !ok {
				return
			}
			s.err <- err
		case <-s.closing:
			return
		}
	}
}

// Service represents a service attached to the server.
type Service interface {
	Open() error
	Close() error
}

// prof stores the file locations of active profiles.
var prof struct {
	cpu *os.File
	mem *os.File
}

// StartProfile initializes the cpu and memory profile, if specified.
func startProfile(cpuprofile, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Errorf("cpuprofile: %s", err)
		}
		log.Info("writing CPU profile to: %s\n", cpuprofile)
		prof.cpu = f
		pprof.StartCPUProfile(prof.cpu)
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Errorf("memprofile: %s", err)
		}
		log.Info("writing mem profile to: %s\n", memprofile)
		prof.mem = f
		runtime.MemProfileRate = 4096
	}

}

// StopProfile closes the cpu and memory profiles if they are running.
func stopProfile() {
	if prof.cpu != nil {
		pprof.StopCPUProfile()
		prof.cpu.Close()
		log.Info("CPU profile stopped")
	}
	if prof.mem != nil {
		pprof.Lookup("heap").WriteTo(prof.mem, 0)
		prof.mem.Close()
		log.Info("mem profile stopped")
	}
}
