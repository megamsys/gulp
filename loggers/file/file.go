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

package file

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/loggers"
	//"github.com/megamsys/gulp/meta"
	"os"
	//"path"
)

func init() {
	loggers.Register("file", fileManager{})
}

type fileManager struct{}

func (m fileManager) Notify(boxName string, messages []loggers.Boxlog, f interface{}) error {
	file := f.(*os.File)
	for _, msg := range messages {
		if _, err := file.WriteString(msg.Message + "\n"); err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
			return err
		}
	}
	return nil
}
