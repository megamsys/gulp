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
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/cmd"
)

type BootOpts struct {
	B *provision.Box
}

func (opts *BootOpts) OK() bool {
	return opts.B.Status == provision.StatusLaunched
}

// Boot runs the boot of the vm.
func Boot(opts *BootOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := NewLogWriter(opts.B)
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)

	err := bootUpBox(opts, writer)
	elapsed := time.Since(start)
	saveErr := saveBootData(opts, outBuffer.String(), elapsed)
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save boot data, boot opts: %#v", opts)
	}

	if err != nil {
		return err
	}
	return nil
}

func bootUpBox(boot *BootOpts, writer io.Writer) error {
	if boot.OK() {
		fmt.Fprintf(writer, "  boot for box (%s)\n", boot.B.GetFullName())
		if bs, ok := Provisioner.(provision.Deployer); ok {
			return bs.Bootstrap(boot.B, writer)
		}
	} else {
		fmt.Fprintf(writer, "  skip boot for box (%s)\n", boot.B.GetFullName())
	}
	fmt.Fprintf(writer, "  boot for box (%s) OK\n", boot.B.GetFullName())
	return nil
}

func saveBootData(boot *BootOpts, blog string, elapsed time.Duration) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(boot.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(blog, "yellow", "", ""))
	return nil
}
