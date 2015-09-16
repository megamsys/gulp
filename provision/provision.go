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

package provision

import (
	"errors"
	"fmt"
	"io"

	"github.com/megamsys/gulp/carton/bind"
)

const (
	PROVIDER        = "provider"
	PROVIDER_CHEF    = "chefsolo"
)

var (
	ErrInvalidStatus = errors.New("invalid status")
	ErrEmptyCarton   = errors.New("no boxs for this carton")
	ErrBoxNotFound   = errors.New("box not found")
)

// Status represents the status of a unit in megamd
type Status string

func (s Status) String() string {
	return string(s)
}

func ParseStatus(status string) (Status, error) {
	switch status {
	case "deploying":
		return StatusDeploying, nil
	case "creating":
		return StatusCreating, nil
	case "error":
		return StatusError, nil
	}
	return Status(""), ErrInvalidStatus
}

const (
	// StatusDeploying is the initial status of a box in the database
	// it should transition shortly to a more specific status
	StatusDeploying = Status("deploying")

	// StatusCreating is the status for box being provisioned by the
	// provisioner, like in the deployment.
	StatusCreating = Status("creating")

	// StatusCreated is the status for box after being provisioned by the
	// provisioner, updated by gulp
	StatusCreated = Status("created")

	// Stateup is the status for box being statefully moved to a different state.
	// Sent by megamd to gulpd when it received StatusCreated.
	StatusStateup = Status("stateup")


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

	Bind(*Box) error
	Unbind(*Box) error

	// Log should be used to log messages in the box.
	Log(message, source, unit string) error

	Boxes() []*Box

	// Run executes the command in box units. Commands executed with this
	// method should have access to environment variables defined in the
	// app.
	Run(cmd string, w io.Writer, once bool) error

	Envs() map[string]bind.EnvVar

	GetMemory() int64
	GetSwap() int64
	GetCpuShare() int
}

// Deployer is a provisioner that can deploy the box from a
type Deployer interface {
    Create(b *Box, w io.Writer) error
	Deploy(b *Box, w io.Writer) error
}


// Provisioner is the basic interface of this package.
//
// Any gulpd provisioner must implement this interface in order to provision
// gulpd cartons.
// A Provisioner is responsible for provisioning a machine with Chef.
type Provisioner interface {
	//PrepareFiles() error
	Command() []string
}

type MessageProvisioner interface {
	StartupMessage() (string, error)
}

// InitializableProvisioner is a provisioner that provides an initialization
// method that should be called when the carton is started,
//additionally provide a map of configuration info.
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
