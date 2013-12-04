package main

import (
	"github.com/indykish/gulp/amqp"
	"log"
	"sync"
)

const queueName = "gulpd-app"

var (
	qfactory amqp.QFactory
	_queue   amqp.Q
	_handler amqp.Handler
	o        sync.Once
)

func RunServer(dry bool) {
	log.Printf("RunServer    = %s", dry)
	handler().Start()
}

func setQueue() {
	var err error
	qfactory, err = amqp.Factory()
	log.Printf("setQueue    = %s", qfactory)

	if err != nil {
		log.Fatalf("Failed to get the queue instance: %s", err)
	}
	_handler, err = qfactory.Handler(handle, queueName)
	if err != nil {
		log.Fatalf("Failed to create the queue handler: %s", err)
	}
	log.Printf("_handler    = %s", _handler)

	_queue, err = qfactory.Get(queueName)
	log.Printf("_queue    = %s", _queue)

	if err != nil {
		log.Fatalf("Failed to get the queue instance: %s", err)
	}
}

func handler() amqp.Handler {
	o.Do(setQueue)
	return _handler
}

// handle is the function called by the queue handler on each message.
func handle(msg *amqp.Message) {
	log.Printf("Hurray I got a message => %s", msg)
	/*	switch msg.Action {
		case RegenerateApprcAndStart:
			fallthrough
		case regenerateApprc:
			if len(msg.Args) < 1 {
				log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
				msg.Delete()
				return
			}
			app, err := ensureAppIsStarted(msg)
			if err != nil {
				log.Print(err)
				return
			}
			msg.Delete()
			app.SerializeEnvVars()
			fallthrough
		case startApp:
			if msg.Action == regenerateApprc {
				break
			}
			if len(msg.Args) < 1 {
				log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
			}
			app, err := ensureAppIsStarted(msg)
			if err != nil {
				log.Print(err)
				return
			}
			err = app.Restart(ioutil.Discard)
			if err != nil {
				log.Printf("Error handling %q. App failed to start:\n%s.", msg.Action, err)
				return
			}
			msg.Delete()
		case BindService:
			err := bindUnit(msg)
			if err != nil {
				log.Print(err)
				return
			}
			msg.Delete()
		default:
			log.Printf("Error handling %q: invalid action.", msg.Action)
			msg.Delete()
		}
	*/

}
