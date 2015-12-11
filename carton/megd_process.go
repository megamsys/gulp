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

package carton

import (
	"bytes"
)

// DeleteProcs represents a command for delete cartons.
type DeleteProcess struct {
	Name string
}

func (s DeleteProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("DELETE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s DeleteProcess) Process(ca *Carton) error {
	if err := ca.Delete(); err != nil {
		return err
	}
	return nil
}

// StartProcs represents a command for starting  cartons.
type StartProcess struct {
	Name string
}

func (s StartProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("START CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s StartProcess) Process(ca *Carton) error {
	if err := ca.LCoperation(START); err != nil {
		return err
	}
	return nil
}

// StopProcs represents a command for stoping  cartons.
type StopProcess struct {
	Name string
}

func (s StopProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STOP CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s StopProcess) Process(ca *Carton) error {
	if err := ca.LCoperation(STOP); err != nil {
		return err
	}
	return nil
}

// RestartProcs represents a command for restarting  cartons.
type RestartProcess struct {
	Name string
}

func (s RestartProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("RESTART CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s RestartProcess) Process(ca *Carton) error {
	if err := ca.LCoperation(RESTART); err != nil {
		return err
	}
	return nil
}

// StateupProcess represents a command for restarting  cartons.
type StateupProcess struct {
	Name string
}

func (s StateupProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STATEUP CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s StateupProcess) Process(ca *Carton) error {
	if err := ca.Stateup(); err != nil {
		return err
	}
	return nil
}

// StatedownProcess represents a command for restarting  cartons.
type StatedownProcess struct {
	Name string
}

func (s StatedownProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("STATEDOWN CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s StatedownProcess) Process(ca *Carton) error {
	if err := ca.Statedown(); err != nil {
		return err
	}
	return nil
}

// CIStateProcess represents a command for continuos integration  cartons.
type CIStateProcess struct {
	Name string
}

func (s CIStateProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("CISTATE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s CIStateProcess) Process(ca *Carton) error {
	if err := ca.CIState(); err != nil {
		return err
	}
	return nil
}
