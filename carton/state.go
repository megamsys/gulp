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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/cmd"
)

type StateOpts struct {
	B     *provision.Box
	Image string
}

// Deploy runs a deployment of an application.
func Stateup(opts *StateOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()

	logWriter := NewLogWriter(opts.B)
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := deployToProvisioner(opts, writer)
	elapsed := time.Since(start)

	saveErr := saveDeployData(opts, outBuffer.String(), elapsed)
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save deploy data, deploy opts: %#v", opts)
	}
	if err != nil {
		return err
	}
	return nil
}

func deployToProvisioner(opts *StateOpts, writer io.Writer) error {
 fmt.Println("_____________  deploy Provisioner  :",opts.B.GetShortTosca())

   switch opts.B.GetShortTosca() {
   case "vertice":
   	if deployer, ok := Provisioner.(provision.Deployer); ok {
   		deployer.Stateup(opts.B, writer)
   	}
   case "bitnami":
   	if deployer, ok := Provisioner.(provision.BitnamiDeployer); ok {
           deployer.StateupBitnami(opts.B, writer)
     }
   }

	return nil
}

func saveDeployData(opts *StateOpts, dlog string, duration time.Duration) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(duration.String(), "green", "", "bold"),
		cmd.Colorfy(dlog, "yellow", "", ""))
	return nil
}
