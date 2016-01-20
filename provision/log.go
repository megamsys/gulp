package provision

import (
	log "github.com/Sirupsen/logrus"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/gulp/meta"
)

const (
	maxInFlight = 150
)

var LogPubSubQueueSuffix = "_log"

func logQueue(boxName string) string {
	return boxName + LogPubSubQueueSuffix
}

func notify(boxName string, messages []interface{}) error {

	pons := nsqp.New()

	if err := pons.Connect(meta.MC.NSQd[0]); err != nil {
		return err
	}

	defer pons.Stop()

	for _, msg := range messages {
		log.Debugf("%s:%s", logQueue(boxName),msg)
		if err := pons.PublishJSONAsync(logQueue(boxName), msg, nil); err != nil {
			log.Errorf("Error on publish: %s", err.Error())
		}
	}
	return nil
}
