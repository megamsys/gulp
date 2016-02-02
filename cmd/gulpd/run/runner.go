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
package run

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	pp "github.com/megamsys/libgo/cmd"
	"github.com/tj/go-spin"
)

const logo = `
 ██████╗ ██╗   ██╗██╗     ██████╗ ██████╗
██╔════╝ ██║   ██║██║     ██╔══██╗██╔══██╗
██║  ███╗██║   ██║██║     ██████╔╝██║  ██║
██║   ██║██║   ██║██║     ██╔═══╝ ██║  ██║
╚██████╔╝╚██████╔╝███████╗██║     ██████╔╝
 ╚═════╝  ╚═════╝ ╚══════╝╚═╝     ╚═════╝
`

// Command represents the command executed by "gulpd start".
type Command struct {
	Version    string
	Branch     string
	Commit     string
	CPUProfile string
	MemProfile string

	closing chan struct{}
	Closed  chan struct{}

	Server *Server
}

func NewCommand() *Command {
	return &Command{
		CPUProfile: "cpuprof",
		MemProfile: "memprof",
		closing:    make(chan struct{}),
		Closed:     make(chan struct{}),
	}
}

func (cmd *Command) Gpd(c *Config, version string) error {
	cmd.funSpin(pp.Colorfy(logo, "green", "", "bold"), version)

	s, err := NewServer(c, cmd.Version)
	if err != nil {
		return fmt.Errorf("create server: %s", err)
	}
	s.CPUProfile = cmd.CPUProfile
	s.MemProfile = cmd.MemProfile
	if err := s.Open(); err != nil {
		return fmt.Errorf("open server: %s", err)
	}
	cmd.Server = s

	go cmd.monitorServerErrors()
	return nil
}

// Close shuts down the server.
func (cmd *Command) Close() error {
	defer close(cmd.Closed)
	close(cmd.closing)
	if cmd.Server != nil {
		return cmd.Server.Close()
	}
	return nil
}

func (cmd *Command) monitorServerErrors() {
	for {
		select {
		case err := <-cmd.Server.Err():
			log.Error(err)
		case <-cmd.closing:
			return
		}
	}
}

func (cmd *Command) funSpin(vers string, logo string) {
	fmt.Printf("%s %s", vers, logo)

	s := spin.New()
	for i := 0; i < 10; i++ {
		fmt.Printf("\r%s", fmt.Sprintf("%s %s", pp.Colorfy("starting", "green", "", "bold"), s.Next()))
		time.Sleep(3 * time.Millisecond)
	}
	fmt.Printf("\n")
}
