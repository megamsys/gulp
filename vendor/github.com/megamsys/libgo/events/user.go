package events

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/events/alerts"
	constants "github.com/megamsys/libgo/utils"
)

// AfterFunc represents a after alert function, that can be registered with
// NewUser function.
type AfterFunc func(evt *Event) error

type AfterFuncs []AfterFunc

type AfterFuncsMap map[alerts.EventAction]AfterFuncs

var notifiers map[string]alerts.Notifier

type User struct {
	stop chan struct{}
	fns  AfterFuncsMap
}

func NewUser(e EventsConfigMap, fnmap AfterFuncsMap) *User {
	register(e)
	return &User{fns: fnmap}
}

func register(e EventsConfigMap) {
	notifiers = make(map[string]alerts.Notifier)
	notifiers[constants.MAILGUN] = newMailgun(e.Get(constants.MAILGUN), e.Get(constants.META))
	notifiers[constants.INFOBIP] = newInfobip(e.Get(constants.INFOBIP))
	notifiers[constants.SLACK] = newSlack(e.Get(constants.SLACK))
	notifiers[constants.SCYLLA] = newScylla(e.Get(constants.META))
}

func newMailgun(m map[string]string, n map[string]string) alerts.Notifier {
	return alerts.NewMailgun(m, n)
}

func newInfobip(m map[string]string) alerts.Notifier {
	return alerts.NewInfobip(m)
}

func newSlack(m map[string]string) alerts.Notifier {
	return alerts.NewSlack(m)
}

func newScylla(m map[string]string) alerts.Notifier {
	return alerts.NewScylla(m)
}

// Watches for new vms, or vms destroyed.
func (self *User) Watch(eventsChannel *EventChannel) error {
	self.stop = make(chan struct{})
	go func() {
		for {
			select {
			case event := <-eventsChannel.channel:
				err := self.alert(event)
				if err != nil {
					log.Warningf("Failed to process watch event: %v", err)
				}
			case <-self.stop:
				log.Info("user watcher exiting")
				return
			}
		}
	}()
	return nil
}

func (self *User) alert(evt *Event) error {
	var err error
	for _, a := range notifiers {
		err = a.Notify(evt.EventAction, evt.EventData)
	}
	if err != nil {
		return err
	}
	return self.after(evt)
}

func (self *User) after(evt *Event) error {
	var err error
	perActionfns := self.fns[evt.EventAction]
	for _, fn := range perActionfns {
		err = fn(evt)
	}
	return err
}

func (self *User) Close() {
	if self.stop != nil {
		close(self.stop)
	}
}
