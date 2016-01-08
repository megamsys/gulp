/*
** copyright [2013-2016] [Megam Systems]
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
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/cmd"
)

type LifecycleOpts struct {
	B         *provision.Box
	start     time.Time
	logWriter LogWriter
	writer    io.Writer
	outBuffer bytes.Buffer
}

func (li *LifecycleOpts) setLogger() {
	li.start = time.Now()
	li.logWriter = NewLogWriter(li.B)
	li.writer = io.MultiWriter(&li.outBuffer, &li.logWriter)
}

//if the state is in running, started, stopped, restarted then allow it to be lcycled.
// to-do: allow states that ends with "*ing or *ed" that should fix this generically.
func (li *LifecycleOpts) canCycle() bool {
	return li.B.Status == provision.StatusRunning ||
		li.B.Status == provision.StatusStarted ||
		li.B.Status == provision.StatusStopped ||
		li.B.Status == provision.StatusStarted ||
		li.B.Status == provision.StatusUpgraded
}

// Starts  the box.
func Start(li *LifecycleOpts) error {
	li.setLogger()
	defer li.logWriter.Close()
	if li.canCycle() {
		if err := Provisioner.Start(li.B, li.writer); err != nil {
			return err
		}
	}
	saveErr := saveLifecycleData(li, li.outBuffer.String(), time.Since(li.start))
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save lifecycle data, lifecycle opts: %#v", li)
	}
	return nil
}

// Stops the box
func Stop(li *LifecycleOpts) error {
	li.setLogger()
	defer li.logWriter.Close()
	if li.canCycle() {
		if err := Provisioner.Stop(li.B, li.writer); err != nil {
			return err
		}
	}
	saveErr := saveLifecycleData(li, li.outBuffer.String(), time.Since(li.start))
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save lifecycle data, lifecycle opts: %#v", li)
	}
	return nil
}

// Restart the box.
func Restart(li *LifecycleOpts) error {
	li.setLogger()
	defer li.logWriter.Close()
	if li.canCycle() {
		if err := Provisioner.Restart(li.B, li.writer); err != nil {
			return err
		}
	}
	saveErr := saveLifecycleData(li, li.outBuffer.String(), time.Since(li.start))
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save lifecycle data, lifecycle opts: %#v", li)
	}
	return nil
}

func saveLifecycleData(li *LifecycleOpts, llog string, elapsed time.Duration) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(li.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(llog, "yellow", "", ""))
	return nil
}
