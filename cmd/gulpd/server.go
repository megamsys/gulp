package main

import (	
	"github.com/indykish/gulp/amqp"
	"log"
	"sync"
)

const queueName = "gulpd-app"


func RunServer(dry bool) {
	q := aqueue()	
	handler().Start()
}


var (
	qfactory amqp.QFactory
	_queue   amqp.Q
	_handler amqp.Handler
	o        sync.Once
)

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

func handler() amqp.Handler {
	o.Do(setQueue)
	return _handler
}

func aqueue() amqp.Q {
	o.Do(setQueue)
	return _queue
}