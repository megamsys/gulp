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

package controls

import (
	"io"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
	)

// Status represents the status of a unit in megamd
type Status string

func (s Status) String() string {
	return string(s)
}

const (

	START = "start"
	STOP = "stop"
	RESTART = "restart"

	// StatusStarted 
	StatusStarted = Status("started")
	
	// StatusStopped 
	StatusStopped = Status("stopped")
	
	// StatusRestarted 
	StatusRestarted = Status("restarted")
	
	// StatusError is the status for units that failed to start, because of
	// a box error.
	StatusError = Status("error")
	
	)


func Restart(box *provision.Box, process string, w io.Writer) error {
	actions := []*action.Action{
		&restart,		
//		&updateStatusInRiak,
	}
	pipeline := action.NewPipeline(actions...)
	
	ctype := strings.Split(box.GetTosca(), ".")
    
	args := runControlActionsArgs{
		box:             box,
		writer:          w,
		machineStatus:   StatusRestarted,
		command: 		 STOP + " " + ctype[2] + "; " + START + " " + ctype[2] + " > /var/log/megam/gulpd.log",
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute create pipeline for box %s - %s", box.GetFullName(), err)
		return err
	}
	return nil	 
}

func Start(box *provision.Box, process string, w io.Writer) error {
	actions := []*action.Action{
		&start,		
//		&updateStatusInRiak,
	}
	pipeline := action.NewPipeline(actions...)
	
	ctype := strings.Split(box.GetTosca(), ".")
    
	args := runControlActionsArgs{
		box:             box,
		writer:          w,
		machineStatus:   StatusStarted,
		command: 		 START + " " + ctype[2] + " > /var/log/megam/gulpd.log",
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute create pipeline for box %s - %s", box.GetFullName(), err)
		return err
	}
	return nil	 
}

func Stop(box *provision.Box, process string, w io.Writer) error {
	actions := []*action.Action{
		&stop,		
//		&updateStatusInRiak,
	}
	pipeline := action.NewPipeline(actions...)
	
	ctype := strings.Split(box.GetTosca(), ".")
    
	args := runControlActionsArgs{
		box:             box,
		writer:          w,
		machineStatus:   StatusStopped,
		command: 		 STOP + " " + ctype[2] + " > /var/log/megam/gulpd.log",
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute create pipeline for box %s - %s", box.GetFullName(), err)
		return err
	}
	return nil	 
}
