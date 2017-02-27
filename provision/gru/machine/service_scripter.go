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

package machine

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/megamsys/libgo/os"
)

type Scriptd struct {
	name    string
	control string
	os      string
}

func NewServiceScripter(name string, control string) *Scriptd {
	return &Scriptd{
		name:    name,
		control: control,
	}
}

func (i *Scriptd) Cmd() []string {
	osh := os.HostOS()
	switch runtime.GOOS {
	case "linux":
		if osh != os.Ubuntu {
			return strings.Fields(fmt.Sprintf("systemctl %s %s", i.control, i.name))
		}
	default:
		return strings.Fields(fmt.Sprintf("%s %s", i.control, i.name))
	}
	return strings.Fields(fmt.Sprintf("%s %s", i.control, i.name))
}
