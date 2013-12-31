package main

import (
	"github.com/indykish/gulp/amqp"
	"github.com/indykish/gulp/app"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

const (
	// queue actions
	startApp   = "start"
	stopApp    = "stop"
	buildApp   = "build"
	restartApp = "restart"
	queueName  = "gulpd-app"
)

var (
	qfactory      amqp.QFactory
	_queue        amqp.Q
	_handler      amqp.Handler
	o             sync.Once
	signalChannel chan<- os.Signal
	nameRegexp    = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
)

func RunServer(dry bool) {
	log.Printf("Gulpd starting at %s", time.Now())
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	handler().Start()
	log.Printf("Gulpd at your service.")
	<-signalChannel
	log.Println("Gulpd killed |_|.")
}

func StopServer(bark bool) {
	log.Printf("Gulpd stopping at %s", time.Now())
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

/*func handle(msg *amqp.Message) {
	log.Printf("Hurray I got a message => %v", msg)

	if nameRegexp.MatchString(msg.Action) {
		msg.Delete()
		manager.Run(msg.Args[1:])
	} else {
		log.Printf("Error handling %q: invalid action.", msg.Action)
	}

}
*/

// handle is the function called by the queue handler on each message.
func handle(msg *amqp.Message) {
	log.Printf("Handling message %v", msg)

	switch msg.Action {
	case restartApp:
		fallthrough
	case startApp:
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Name: "myapp", Id: "RIPAB"}

		if err := ap.Get(msg.Id); err != nil {
			log.Printf("Error handling %q: Riak didn't cooperate:\n%s.", msg.Action, err)
			return
		}
        log.Printf("Handling message %v", ap.GetAppReqs())
		err := app.StartApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to start:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break
	case stopApp:
	    if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Name: "myapp", Id: "RIPAB"}

		if err := ap.Get(msg.Id); err != nil {
			log.Printf("Error handling %q: Riak didn't cooperate:\n%s.", msg.Action, err)
			return
		}
        log.Printf("Handling message %v", ap.GetAppReqs())
		err := app.StopApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to stop:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break
	/*	err := bindUnit(msg)
		if err != nil {
			log.Print(err)
			return
		}
		msg.Delete()
		break
	*/
	default:
		log.Printf("Error handling %q: invalid action.", msg.Action)
		msg.Delete()
	}
}
