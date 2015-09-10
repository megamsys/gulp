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
package activities

import (
	"fmt"
	"io/ioutil"
	libhttp "net/http"

	"github.com/megamsys/gulp/app"
)

// Every Activities must implement this interface.
type Activities interface {
	// Called when execute action a Machine.
	Action(*app.ActionData) error
}

var plugs = make(map[string]Activities)
var plug_names = []string{"state", "control", "policies", "docker"}

/**
**register the all Activities to "plug" array
**/
func RegisterActivities(name string, activity Activities) {
	plugs[name] = activity
}

func GetActivity(name string) (Activities, error) {
	activity, ok := plugs[name]
	if !ok {
		return nil, fmt.Errorf("Activities not registered")
	}
	return activity, nil
}

func HttpDataParser(w libhttp.ResponseWriter, req *libhttp.Request) []byte {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("error")
	}

	return body
}
