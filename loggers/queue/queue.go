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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/loggers"
	"github.com/megamsys/gulp/meta"
	libamqp "github.com/megamsys/libgo/amqp"
	amqp "github.com/streadway/amqp"
)

var LogPubSubQueueSuffix = "_log"

func logQueue(boxName string) string {
	return boxName + LogPubSubQueueSuffix
}

func init() {
	loggers.Register("amqp",
		func() interface{} {
			return new(AMQPOutput)
		})
}

// AMQP Output config struct
type AMQPOutputConfig struct {
	// AMQP URL. Spec: http://www.rabbitmq.com/uri-spec.html
	// Ex: amqp://USERNAME:PASSWORD@HOSTNAME:PORT/
	URL string
	// Exchange name
	Exchange string
	// Type of exchange, options are: fanout, direct, topic, headers
	ExchangeType string
	// Whether the exchange should be durable or not
	// Defaults to non-durable
	ExchangeDurability bool
	// Whether the exchange is deleted when all queues have finished
	// Defaults to auto-delete
	ExchangeAutoDelete bool
	// Routing key for the message to send, or when used for consumer
	// the routing key to bind the queue to the exchange with
	// Defaults to empty string
	RoutingKey string
	// Whether messages published should be marked as persistent or
	// transient. Defaults to non-persistent.
	Persistent bool
	// Whether messages published should be fully serialized when
	// published. The AMQP input will automatically detect these
	// messages and deserialize them. Defaults to true.
	Serialize bool
	// Optional subsection for TLS configuration of AMQPS connections. If
	// unspecified, the default AMQPS settings will be used.
	Tls tcp.TlsConfig
	// MIME content type for the AMQP header.
	ContentType string
	// Allows us to use framing by default.
	UseFraming bool
}

type AMQPOutput struct {
	// Hold a copy of the config used
	config *AMQPOutputConfig
	// The AMQP Channel created upon Init
	ch AMQPChannel
	// closeChan gets sent an error should the channel close so that
	// we can immediately exit the output
	closeChan chan *amqp.Error
	usageWg   *sync.WaitGroup
	// connWg tracks whether the connection is no longer in use
	// and is used as a barrier to ensure all users of the connection
	// are done before we finish
	connWg *sync.WaitGroup
	// Hold a reference to the connection hub.
	amqpHub libamqp.AMQPConnectionHub
}

func (ao *AMQPOutput) ConfigStruct() interface{} {
	return &AMQPOutputConfig{
		ExchangeDurability: false,
		ExchangeAutoDelete: true,
		RoutingKey:         "",
		Persistent:         false,
		Encoder:            "ProtobufEncoder",
		UseFraming:         true,
		ContentType:        "application/json",
	}
}

func (ao *AMQPOutput) Init(config interface{}) (err error) {
	conf := config.(*AMQPOutputConfig)
	ao.config = conf
	var tlsConf *tls.Config = nil
	if strings.HasPrefix(conf.URL, "amqps://") && &ao.config.Tls != nil {
		if tlsConf, err = tcp.CreateGoTlsConfig(&ao.config.Tls); err != nil {
			return fmt.Errorf("TLS init error: %s", err)
		}
	}

	var dialer = NewAMQPDialer(tlsConf)
	if ao.amqpHub == nil {
		ao.amqpHub = GetAmqpHub()
	}
	ch, usageWg, connectionWg, err := ao.amqpHub.GetChannel(conf.URL, dialer)
	if err != nil {
		return
	}
	ao.connWg = connectionWg
	ao.usageWg = usageWg
	closeChan := make(chan *amqp.Error)
	ao.closeChan = ch.NotifyClose(closeChan)
	err = ch.ExchangeDeclare(conf.Exchange, conf.ExchangeType,
		conf.ExchangeDurability, conf.ExchangeAutoDelete, false, false,
		nil)
	if err != nil {
		usageWg.Done()
		return
	}
	ao.ch = ch
	return
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

func (ao *AMQPOutput) Run(or OutputRunner, h PluginHelper) (err error) {
	if or.Encoder() == nil {
		return errors.New("Encoder required.")
	}

	inChan := or.InChan()
	conf := ao.config

	var (
		pack     *PipelinePack
		persist  uint8
		ok       bool = true
		amqpMsg  amqp.Publishing
		outBytes []byte
	)
	if conf.Persistent {
		persist = amqp.Persistent
	} else {
		persist = amqp.Transient
	}

	// Spin up separate goroutine so we can wait for close notifications from
	// the AMQP lib w/o deadlocking on our `AMQPChannel.Publish` call.
	stopChan := make(chan struct{})
	go func() {
		<-ao.closeChan
		close(stopChan)
	}()

	for ok {
		select {
		case <-stopChan:
			ok = false
		case pack, ok = <-inChan:
			if !ok {
				break
			}
			if outBytes, err = or.Encode(pack); err != nil {
				or.UpdateCursor(pack.QueueCursor)
				err = fmt.Errorf("Error encoding message: %s", err)
				pack.Recycle(err)
				continue
			} else if outBytes == nil {
				or.UpdateCursor(pack.QueueCursor)
				pack.Recycle(nil)
				continue
			}
			amqpMsg = amqp.Publishing{
				DeliveryMode: persist,
				Timestamp:    time.Now(),
				ContentType:  conf.ContentType,
				Body:         outBytes,
			}
			err = ao.ch.Publish(conf.Exchange, conf.RoutingKey,
				false, false, amqpMsg)
			if err != nil {
				err = NewRetryMessageError(err.Error())
				ok = false
			} else {
				or.UpdateCursor(pack.QueueCursor)
			}
			pack.Recycle(err)
		}
	}
	ao.usageWg.Done()
	ao.amqpHub.Close(conf.URL, ao.connWg)
	ao.connWg.Wait()
	<-stopChan
	return
}

func (ao *AMQPOutput) CleanupForRestart() {
	return
}
