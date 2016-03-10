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

package provision

import (
	"errors"
	"fmt"
	"io"

	"github.com/megamsys/libgo/exec"
)

const (
	PROVIDER = "provider"
	CHEFSOLO = "chefsolo"
)

var (
	ErrEmptyCarton    = errors.New("no boxs for this carton")
	ErrBoxNotFound    = errors.New("box not found")
	ErrNoOutputsFound = errors.New("no outputs found in the box. Did you set it ?")
)

// Status represents the status of a unit in megamd
type Status string

func (s Status) String() string {
	return string(s)
}

const (
	// StatusLaunched is the status for box being bootted the first time.
	StatusLaunched = Status("launched")

	// StatusBootstrapping is the status for box being bootstrapping by the
	// provisioner, like in the bootstrap.
	StatusBootstrapping = Status("bootstrapping")
	// StatusBootstrapped is the status for box after being bootstrapped by the
	// provisioner, updated by gulp
	StatusBootstrapped = Status("bootstrapped")

	// Stateup is the status for box being statefully moved to a different state.
	// Sent by megamd to gulpd when it received StatusCreated.
	StatusStateup = Status("stateup")

	// StatusRunning
	StatusRunning = Status("running")

	StatusStarting = Status("starting")
	StatusStarted  = Status("started")

	StatusStopping = Status("stopping")
	StatusStopped  = Status("stopped")

	StatusRestarting = Status("restarting")
	StatusRestarted  = Status("restarted")

	StatusUpgrading = Status("upgrading")
	StatusUpgraded  = Status("upgraded")

	// StatusError is the status for units that failed to start, because of
	// a box error.
	StatusError = Status("error")
)

// Named is something that has a name, providing the GetName method.
type Named interface {
	GetName() string
}

// Carton represents a deployment entity in gulpd.
//
// It contains boxes to provision and only relevant information for provisioning.
type Carton interface {
	Named

	Boot() error
	Stateup() error
	Start() error
	Stop() error
	Upgrade() error

	// Run executes the command in box units. Commands executed with this
	// method should have access to environment variables defined in the
	// app.
	Run(cmd string, w io.Writer, once bool) error
}

// Deployer is a provisioner that can deploy the box from a
type Deployer interface {
	Bootstrap(b *Box, w io.Writer) error
	Stateup(b *Box, w io.Writer) error
}

// Provisioner is the basic interface of this package.
//
// A Provisioner is responsible for managing the state of the machine.
type Provisioner interface {
	Command() []string
	Start(b *Box, w io.Writer) error
	Stop(b *Box, w io.Writer) error
	Restart(b *Box, w io.Writer) error
}

// Provisioner message
type MessageProvisioner interface {
	StartupMessage() (string, error)
}

// InitializableProvisioner is a provisioner that provides an initialization
// method that should be called when the carton is started,
// additionally provide a map of configuration info.
type InitializableProvisioner interface {
	Initialize(m map[string]string) error
}

var provisioners = make(map[string]Provisioner)

// Register registers a new provisioner in the Provisioner registry.
func Register(name string, p Provisioner) {
	provisioners[name] = p
}

// Get gets the named provisioner from the registry.
func Get(name string) (Provisioner, error) {
	p, ok := provisioners[name]
	if !ok {
		return nil, fmt.Errorf("unknown provisioner: %q", name)
	}
	return p, nil
}

// Registry returns the list of registered provisioners.
func Registry() []Provisioner {
	registry := make([]Provisioner, 0, len(provisioners))
	for _, p := range provisioners {
		registry = append(registry, p)
	}
	return registry
}

// Error represents a provisioning error. It encapsulates further errors.
type Error struct {
	Reason string
	Err    error
}

// Error is the string representation of a provisioning error.
func (e *Error) Error() string {
	var err string
	if e.Err != nil {
		err = e.Err.Error() + ": " + e.Reason
	} else {
		err = e.Reason
	}
	return err
}

func ExecuteCommandOnce(commandWords []string, w io.Writer) error {
	var e exec.OsExecutor

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, w, w); err != nil {
			return err
		}
	}
	return nil
}
