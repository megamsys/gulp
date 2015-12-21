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

	"github.com/megamsys/gulp/loggers/file"
	"github.com/megamsys/gulp/loggers/queue"
	"github.com/megamsys/gulp/provision"
//	"github.com/megamsys/gulp/operations"
)

type DeployOpts struct {
	B     *provision.Box
	Image string
}

// Deploy runs a deployment of an application.
func Deploy(opts *DeployOpts) error {
	var outBuffer bytes.Buffer
	//	start := time.Now()

	queueWriter := queue.LogWriter{Box: opts.B, Source: opts.B.GetFullName()}
	queueWriter.Async()
	defer queueWriter.Close()

	fileWriter := file.LogWriter{Box: opts.B, Source: opts.B.GetFullName()}
	fileWriter.Async()
	defer fileWriter.Close()

	writer := io.MultiWriter(&outBuffer, &queueWriter, &fileWriter)
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

func BindService(opts *DeployOpts) error {
	/*a, err := operations.Get(args.operations)
	if err != nil {
	//	log.Errorf("fatal error, couldn't locate the Repository %s", args.operations)
		return err
	}
	if binder, ok := operatins.Operation.(opertations.Binder); ok {
	//if binder, ok := Operations.(opertations.Binder); ok {
		return binder.Apply()
	}*/
return nil
}
