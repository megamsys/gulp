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

package queue

import (
	//"fmt"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/loggers"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/libgo/amqp"
	//"strings"
)

func init() {
	loggers.Register("queue", queueManager{})
}

type queueManager struct{}

var LogPubSubQueueSuffix = "_log"

func logQueue(boxName string) string {
	return boxName + LogPubSubQueueSuffix
}

func (m queueManager) Notify(boxName string, messages []loggers.Boxlog) error {
	pubSubQ, err := amqp.NewRabbitMQ(meta.MC.AMQP, logQueue(boxName))
	if err != nil {
		return err
	}

	for _, msg := range messages {
		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
			continue
		}
		err = pubSubQ.Pub(bytes)
		if err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
		}
	}
	return nil
}
