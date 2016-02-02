package provision

import (
	log "github.com/Sirupsen/logrus"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/gulp/meta"
	"unicode/utf8"
)

const (
	maxInFlight = 300
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
		log.Debugf("%s:%s", logQueue(boxName), msg)
		if err := pons.PublishJSONAsync(logQueue(boxName), utftostring(msg.(string)), nil); err != nil {
			log.Errorf("Error on publish: %s", err.Error())
		}
	}
	return nil
}

func utftostring(s string) string {
	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for i, r := range s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		s = string(v)
	}
	return s
}
