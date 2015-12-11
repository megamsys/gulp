// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package file

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/meta"
	"os"
	"path"
	"time"
)

var f *os.File

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

	basePath := meta.MC.Dir + "/logs"
	dir := path.Join(basePath, w.Source)

	filePath := path.Join(dir, w.Source+"_log")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Debugf("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return
		}
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Errorf("Error on logs notify: %s", err.Error())
		return
	}

	go func(f *os.File) {
		defer close(w.doneCh)
		defer f.Close()
		for msg := range w.msgCh {
			err := w.write(msg, f)
			if err != nil {
				log.Errorf("[LogWriter] failed to write async logs: %s", err)
				return
			}
		}
	}(f)
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
		return len(data), w.write(data, f)
	}
	copied := make([]byte, len(data))
	copy(copied, data)
	w.msgCh <- copied
	return len(data), nil
}

func (w *LogWriter) write(data []byte, f *os.File) error {
	return w.Box.Log(string(data), "file", "box", f)
}
