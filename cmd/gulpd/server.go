package main

import (
	"github.com/indykish/gulp/amqp"
	"time"
	"os"
	"os/signal"
	"syscall"
	"log"
	"sync"
)

const queueName = "gulpd-app"

var (
	qfactory amqp.QFactory
	_queue   amqp.Q
	_handler amqp.Handler
	o        sync.Once
	signalChannel chan<- os.Signal
)

func RunServer(dry bool) {
	log.Printf("Gulpd starting at %s",time.Now())
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	handler().Start()
	log.Printf("Gulpd at your service.")
	<-signalChannel
	log.Println("Gulpd killed |_|.")
}

func StopServer(bark bool) {
	log.Printf("Gulpd stopping at %s",time.Now())
	handler().Stop()
	close(signalChannel)
	log.Println("Gulpd finished |-|.")
}


func setQueue() {
	var err error
	qfactory, err = amqp.Factory()

	if err != nil {
		log.Fatalf("Failed to get the queue instance: %s", err)
	}
	_handler, err = qfactory.Handler(handle, queueName)
	if err != nil {
		log.Fatalf("Failed to create the queue handler: %s", err)
	}

	_queue, err = qfactory.Get(queueName)

	if err != nil {
		log.Fatalf("Failed to get the queue instance: %s", err)
	}
}

func aqueue() amqp.Q {
	o.Do(setQueue)
	return _queue
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
