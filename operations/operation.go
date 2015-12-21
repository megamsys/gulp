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

package operations

import (
	"fmt"
)

type JsonPairs []*JsonPair

type JsonPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

func NewJsonPair(k string, v string) *JsonPair {
	return &JsonPair{
		K: k,
		V: v,
	}
}

//match for a value in the JSONPair and send the value
func (p *JsonPairs) match(k string) string {
	for _, j := range *p {
		if j.K == k {
			return j.V
		}
	}
	return ""
}

var managers map[string]InitializableOperation

/* Operation represents a various operations of the application. */
type Operate struct {
	OperationType         string    `json:"operation_type"`
	Description           string    `json:"description"`
	OperationRequirements JsonPairs `json:"operation_requirements"`
}

func (o Operate) GetType() string {
	return o.OperationType
}

func (o Operate) GetDescription() string {
	return o.Description
}

type Operation interface {
}

type Binder interface {
Apply() (string, error)
}

type InitializableOperation interface {
	Initialize(operationtype string) error
}

// Get gets the named provisioner from the registry.
func Get(name string) (Operation, error) {
	p, ok := managers[name]
	if !ok {
		return nil, fmt.Errorf("unknown operation: %q", name)
	}
	return p, nil
}

// Manager returns the current configured manager, as defined in the
// configuration file.
func Manager(managerName string) InitializableOperation {
	if _, ok := managers[managerName]; !ok {
		managerName = "nop"
	}
	return managers[managerName]
}

// Register registers a new operation manager, that can be later configured
// and used.
func Register(name string, manager InitializableOperation) {
	if managers == nil {
		managers = make(map[string]InitializableOperation)
	}
	managers[name] = manager
}
