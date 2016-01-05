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

package controls

import (
	"fmt"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
)

const (
	CONTROL = "control"

	// start control is represents lifecycle start operation of virtual machine
	START = "start"

	// stop control is represents lifecycle stop operation of virtual machine
	STOP = "stop"

	// restart control is represents lifecycle restart operation of virtual machine
	RESTART = "restart"
)

func ParseControl(box *provision.Box, action string, w io.Writer) (provision.Status, error) {
	switch action {
	case START:
		return cstart(box, w)
	case STOP:
		return cstop(box, w)
	case RESTART:
		return crestart(box, w)
	default:
		return "", newParseError([]string{CONTROL, action}, []string{START, STOP, RESTART})
	}
}

func crestart(box *provision.Box, w io.Writer) (provision.Status, error) {
	actions := []*action.Action{
		&restart,
	}
	pipeline := action.NewPipeline(actions...)

	ctype := strings.Split(box.GetTosca(), ".")

	args := runControlActionsArgs{
		box:     box,
		writer:  w,
		command: STOP + " " + ctype[2] + "; " + START + " " + ctype[2],
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute create pipeline for box %s - %s", box.GetFullName(), err)
		return "", err
	}
	return provision.StatusRestarted, nil
}

func cstart(box *provision.Box, w io.Writer) (provision.Status, error) {
	actions := []*action.Action{
		&start,
	}
	pipeline := action.NewPipeline(actions...)

	ctype := strings.Split(box.GetTosca(), ".")

	args := runControlActionsArgs{
		box:     box,
		writer:  w,
		command: START + " " + ctype[2],
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute create pipeline for box %s - %s", box.GetFullName(), err)
		return "", err
	}
	return provision.StatusStarted, nil
}

func cstop(box *provision.Box, w io.Writer) (provision.Status, error) {
	actions := []*action.Action{
		&stop,
	}
	pipeline := action.NewPipeline(actions...)

	ctype := strings.Split(box.GetTosca(), ".")

	args := runControlActionsArgs{
		box:     box,
		writer:  w,
		command: STOP + " " + ctype[2],
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("error on execute create pipeline for box %s - %s", box.GetFullName(), err)
		return "", err
	}
	return provision.StatusStopped, nil
}

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Found    string
	Expected []string
}

// newParseError returns a new instance of ParseError.
func newParseError(found []string, expected []string) *ParseError {
	return &ParseError{Found: strings.Join(found, ","), Expected: expected}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	return fmt.Sprintf("found %s, expected %s", e.Found, strings.Join(e.Expected, ", "))
}
