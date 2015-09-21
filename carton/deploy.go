/*
** copyright [2013-2015] [Megam Systems]
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
//	"time"
	"fmt"
//	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
)

type DeployOpts struct {
	B      *provision.Box
	Image 	string
}

// Deploy runs a deployment of an application. 
func Deploy(opts *DeployOpts) error {
	var outBuffer bytes.Buffer
//	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	fmt.Println("-------------------------------")
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := deployToProvisioner(opts, writer)
	
	if err != nil {
		return err
	}
	return nil
}

func deployToProvisioner(opts *DeployOpts, writer io.Writer) error {
	if deployer, ok := Provisioner.(provision.Deployer); ok {
		return deployer.Deploy(opts.B, writer)
	}
	return nil
}



