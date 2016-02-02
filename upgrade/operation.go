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

package upgrade

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/gulp/repository"
)

// ErrDuplicateOperation is the error returned by Register when the given name
// is already in use.
var ErrDuplicateOperation = errors.New("there's already a operation with this name")

//duplicate flags as they are in provision.StatusUpgraded as well
const StatusUpgraded = "upgraded"
const StatusError = "error"

type Operation struct {
	Type        string         `json:"operation_type"`
	Description string         `json:"description"`
	Properties  bind.JsonPairs `json:"properties"`
	Status      string         `json:"status"`
}

type operation struct {
	Name     string
	Ran      bool
	Optional bool
	fn       OperateFunc
}

// OperateFunc represents a operation function, that can be registered with the
// Register function. Operations are later ran in the registration order, and
// this package keeps track of which ate have ran already.
type OperateFunc func() error

var operations []operation

// RunArgs is used by Run and RunOptional functions to modify how operations
// are executed.
type RunArgs struct {
	Id     string
	Name   string
	O      []*Operation
	Writer io.Writer
	Force  bool
}

type OperationsToRun struct {
	Raw      *Operation
	Ran      bool
	Optional bool
	fn       OperateFunc
}

type OperationsRan []OperationsToRun

func (o *Operation) Ran() bool {
	return o.Status == StatusUpgraded
}

func (o *Operation) prepBuildHook() *repository.Hook {
	return &repository.Hook{
		Enabled:  true,
		Token:    o.Properties.Match(repository.TOKEN),
		UserName: o.Properties.Match(repository.USERNAME),
	}
}

func (op OperationsRan) Successful() bool {
	var success = false

	if len(op) > 0 {
		success = true
	}
	for _, pastOpsRun := range op {
		success = success && (pastOpsRun.Raw.Status == StatusUpgraded)
	}
	return success
}

func BuildHook(ops []*Operation, opsType string) *repository.Hook {
	for _, o := range ops {
		switch o.Type {
		case opsType:
			return o.prepBuildHook()
		}
	}
	return nil
}

// Register register a new operation for later execution with the Run
// functions.
func Register(name string, fn OperateFunc) error {
	return register(name, false, fn)
}

func register(name string, optional bool, fn OperateFunc) error {
	for _, m := range operations {
		if m.Name == name {
			return nil
		}
	}
	operations = append(operations, operation{Name: name, Optional: optional, fn: fn})
	return nil
}

// Run runs all registered non optional operations
func Run(args RunArgs) (OperationsRan, error) {
	return run(args)
}

func run(args RunArgs) (OperationsRan, error) {
	operationsToRun, err := getOperations(args.O, true)
	if err != nil {
		return nil, err
	}
	for _, m := range operationsToRun {
		if m.Optional {
			continue
		}
		if !m.Ran || !args.Force {
			fmt.Fprintf(args.Writer, "Running operation (%q)...\n", m.Raw.Type)
			err := m.fn()
			if err != nil {
				m.Raw.Status = StatusError
				return nil, err
			}
			m.Ran = true
			m.Raw.Status = StatusUpgraded
			fmt.Fprintf(args.Writer, "Ran operation (%s) OK\n", m.Raw.Type)
		} else {
			fmt.Fprintf(args.Writer, "Skip operation (%s) OK\n", m.Raw.Type)
		}
	}
	return operationsToRun, nil
}

func getOperations(ran []*Operation, ignoreRan bool) ([]OperationsToRun, error) {
	result := make([]OperationsToRun, 0, len(ran))
	for _, m := range operations {
		m.Ran = false
		for _, r := range ran {
			if strings.EqualFold(r.Type, m.Name) {
				if !ignoreRan || !m.Ran {
					opr := OperationsToRun{
						Raw:      r,
						Ran:      r.Ran(),
						fn:       m.fn,
						Optional: m.Optional,
					}
					result = append(result, opr)
				}
			}
		}
	}
	return result, nil
}
