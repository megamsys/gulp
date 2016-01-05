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

	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/gulp/repository"
)

type Operate struct {
	Type        string         `json:"operation_type"`
	Description string         `json:"description"`
	Properties  bind.JsonPairs `json:"properties"`
}

func (o Operate) GetType() string {
	return o.Type
}

func (o Operate) GetDescription() string {
	return o.Description
}

type Operation interface {
}

// ErrDuplicateOperation is the error returned by Register when the given name
// is already in use.
var ErrDuplicateOperation = errors.New("there's already a operation with this name")

// ErrOperationNotFound is the error returned by RunOptional when the given
// name is not a registered operation.
var ErrOperationNotFound = errors.New("operation not found")

// ErrOperationMandatory is the error returned by Run when the given name is
// not an optional operation. It should be executed calling Run.
var ErrOperationMandatory = errors.New("operation is mandatory")

// ErrOperationAlreadyExecuted is the error returned by Run when the given
// name was previously executed and the force parameter was not supplied.
var ErrOperationAlreadyExecuted = errors.New("operation already executed")

// ErrCannotForceMandatory is the error returned by Run when the force
// paramter is supplied without the name of a operation to run.
var ErrCannotForceMandatory = errors.New("mandatory operations can only run once")

// OperateFunc represents a operation function, that can be registered with the
// Register function. Operations are later ran in the registration order, and
// this package keeps track of which ate have ran already.
type OperateFunc func() error



func BuildHook(ops []*Operate, opsType string) *repository.Hook {
	for _, o := range ops {
		switch o.Type {
		case opsType:
			return o.prepBuildHook()
		}
	}
	return nil
}

func (o *Operate) prepBuildHook() *repository.Hook {
	return &repository.Hook{
		Enabled:  true,
		Token:    o.Properties.Match(repository.TOKEN),
		UserName: o.Properties.Match(repository.USERNAME),
	}
}

// RunArgs is used by Run and RunOptional functions to modify how operations
// are executed.
type RunArgs struct {
	Name   string
	Writer io.Writer
	Force  bool
}

type operation struct {
	Name     string
	Ran      bool
	Optional bool
	fn       UpgradeFunc
}

var operations []operation

// Register register a new operation for later execution with the Run
// functions.
func Register(name string, fn OperateFunc) error {
	return register(name, false, fn)
}

// RegisterOptional register a new operation that will not run automatically
// when calling the Run funcition.
func RegisterOptional(name string, fn OperateFunc) error {
	return register(name, true, fn)
}

func register(name string, optional bool, fn OperateFunc) error {
	for _, m := range operations {
		if m.Name == name {
			return ErrDuplicateOperation
		}
	}
	operations = append(operations, operation{Name: name, Optional: optional, fn: fn})
	return nil
}

// Run runs all registered non optional operations if no ".Name" is informed.
// Migrations are executed in the order that they were registered. If ".Name"
// is informed, an optional operation with the given name is executed.
func Run(args RunArgs) error {
	if args.Name != "" {
		return runOptional(args)
	}
	if args.Force {
		return ErrCannotForceMandatory
	}
	return run(args)
}

func run(args RunArgs) error {
	operationsToRun, err := getOperations(true)
	if err != nil {
		return err
	}
	coll, err := collection()
	if err != nil {
		return err
	}
	defer coll.Close()
	for _, m := range operationsToRun {
		if m.Optional {
			continue
		}
		fmt.Fprintf(args.Writer, "Running %q... ", m.Name)
		err := m.fn()
		if err != nil {
			return err
		}
		m.Ran = true
		err = coll.Insert(m)
		if err != nil {
			return err
		}
		fmt.Fprintln(args.Writer, "OK")
	}
	return nil
}

func runOptional(args RunArgs) error {
	operationsToRun, err := getOperations(false)
	if err != nil {
		return err
	}
	var toRun *operation
	for i, m := range operationsToRun {
		if m.Name == args.Name {
			toRun = &operationsToRun[i]
			break
		}
	}
	if toRun == nil {
		return ErrOperationNotFound
	}
	if !toRun.Optional {
		return ErrOperationMandatory
	}
	if toRun.Ran && !args.Force {
		return ErrOperationAlreadyExecuted
	}
	fmt.Fprintf(args.Writer, "Running %q... ", toRun.Name)
	coll, err := collection()
	if err != nil {
		return err
	}
	defer coll.Close()
	err = toRun.fn()
	if err != nil {
		return err
	}
	toRun.Ran = true
	_, err = coll.Upsert(bson.M{"name": toRun.Name}, toRun)
	if err != nil {
		return err
	}
	fmt.Fprintln(args.Writer, "OK")
	return nil
}

func List() ([]operation, error) {
	return getOperations(false)
}

func getOperations(ignoreRan bool) ([]operation, error) {
	coll, err := collection()
	if err != nil {
		return nil, err
	}
	defer coll.Close()
	result := make([]operation, 0, len(operations))
	names := make([]string, len(operations))
	for i, m := range operations {
		names[i] = m.Name
	}
	query := bson.M{"name": bson.M{"$in": names}, "ran": true}
	var ran []operation
	err = coll.Find(query).All(&ran)
	if err != nil {
		return nil, err
	}
	for _, m := range operations {
		m.Ran = false
		for _, r := range ran {
			if r.Name == m.Name {
				m.Ran = true
				break
			}
		}
		if !ignoreRan || !m.Ran {
			result = append(result, m)
		}
	}
	return result, nil
}

func collection() (*storage.Collection, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	return conn.Collection("operations"), nil
}
