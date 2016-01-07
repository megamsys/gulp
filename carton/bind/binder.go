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
package bind

import (
	"fmt"
	"runtime"

	"github.com/megamsys/libgo/os"
)

// EnvVar represents a environment variable for a carton.
type EnvVar struct {
	Name     string
	Value    string
	Endpoint string
}

func (e *EnvVar) String() string {
	return fmt.Sprintf("%s=%s", e.Name, e.Value)
}

type EnvVars  []EnvVar

func (en EnvVars) WrapForInitds() string {
	var envs = ""
	for _,de := range en {
		envs += wrapForInitdservice(de.Name, de.Value)
	}
	return envs
}

func wrapForInitdservice(key string, value string) string {
	osh := os.HostOS()
	switch runtime.GOOS {
	case "linux":
		if osh != os.Ubuntu {
			return key + "=" + value //systemd
		}
	default:
		return "initctl set-env " + key + "=" + value + "\n"
	}
	return "initctl set-env " + key + "=" + value + "\n"
}


type JsonPair struct {
	K string `json:"key"`
	V string `json:"value"`
}

type JsonPairs []*JsonPair

func NewJsonPair(k string, v string) *JsonPair {
	return &JsonPair{
		K: k,
		V: v,
	}
}

//match for a value in the JSONPair and send the value
func (p *JsonPairs) Match(k string) string {
	for _, j := range *p {
		if j.K == k {
			return j.V
		}
	}
	return ""
}

//match for a value in the JSONPair and send the value
func (p *JsonPairs) NukeAndSet(m map[string][]string) {
	for mkey, mvals := range m { //r is key, s is value
		for i, j := range *p { //i index, j is value
			if j.K == mkey {
				p1 := *p
				p1 = append(p1[:i], p1[i+1:]...) //nuke the old ones that match r.
			}
		}
		for _, mval := range mvals {
			*p = append(*p, NewJsonPair(mkey, mval))
		}
	}
}
