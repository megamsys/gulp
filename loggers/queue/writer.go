// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/libgo/amqp"
	"time"
)

var pub amqp.PubSubQ

type Logger interface {
	Log(string, string, string, interface{}) error
}

type LogWriter struct {
	Box    Logger
	Source string
	msgCh  chan []byte
	doneCh chan bool
}

func (w *LogWriter) Async() {
	w.msgCh = make(chan []byte, 1000)
	w.doneCh = make(chan bool)

	pub, err := amqp.NewRabbitMQ(meta.MC.AMQP, logQueue(w.Source))
	if err != nil {
		return
	}
	err = pub.Connect()
	if err != nil {
		return
	}

	go func(pub amqp.PubSubQ) {
		defer close(w.doneCh)
		defer pub.Close()
		for msg := range w.msgCh {
			err := w.write(msg, pub)
			if err != nil {
				log.Errorf("[LogWriter] failed to write async logs: %s", err)
				return
			}
		}
	}(pub)
}

func (w *LogWriter) Close() {
	if w.msgCh != nil {
		close(w.msgCh)
	}
}

func (w *LogWriter) Wait(timeout time.Duration) error {
	if w.msgCh == nil {
		return nil
	}
	select {
	case <-w.doneCh:
	case <-time.After(timeout):
		return errors.New("timeout waiting for writer to finish")
	}
	return nil
}

// Write writes and logs the data.
func (w *LogWriter) Write(data []byte) (int, error) {
	if w.msgCh == nil {
		return len(data), w.write(data, pub)
	}
	copied := make([]byte, len(data))
	copy(copied, data)
	w.msgCh <- copied
	return len(data), nil
}

func (w *LogWriter) write(data []byte, pub amqp.PubSubQ) error {
	return w.Box.Log(string(data), "queue", "box", pub)
}
