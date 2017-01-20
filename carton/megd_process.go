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

package carton

import (
	"bytes"
)

// BootProcess represents the initial boot for applying cartons.
type BootProcess struct {
	Name string
}

func (s BootProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("BOOT CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s BootProcess) Process(ca *Carton) error {
	if err := ca.Boot(); err != nil {
		return err
	}
	return nil
}

// DeleteProcs represents a command for delete cartons.
type DeleteProcess struct {
	Name string
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
	if err := ca.Start(); err != nil {
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
	if err := ca.Stop(); err != nil {
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
	if err := ca.Restart(); err != nil {
		return err
	}
	return nil
}

// UpgradeProcs represents a command for starting  cartons.
type UpgradeProcess struct {
	Name string
}

func (s UpgradeProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("UPGRADE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s UpgradeProcess) Process(ca *Carton) error {
	if err := ca.Upgrade(); err != nil {
		return err
	}
	return nil
}


// ResetPassword to reset new password of VM root user.
type ResetPasswordProcess struct {
	Name string
}

func (s ResetPasswordProcess) String() string {
	var buf bytes.Buffer
	_, _ = buf.WriteString("UPGRADE CARTON ")
	_, _ = buf.WriteString(s.Name)
	return buf.String()
}

func (s ResetPasswordProcess) Process(ca *Carton) error {
	if err := ca.ResetPassword(); err != nil {
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
