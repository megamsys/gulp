/*
** Copyright [2012-2013] [Megam Systems]
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
package amqp

import (
	"errors"
	"fmt"
	"github.com/globocom/config"
	"github.com/streadway/amqp"
	"log"
	"net"
	"regexp"
	"sync"
	"time"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
	tag     string
	done    chan error
}

type rabbitmqQ struct {
	name string
}

const (
	DefaultAMQPURL      = "http://localhost:8098/amqp"
	DefaultQueue        = "thirsting.megam.co"
	DefaultExchange     = "megam_nodes"
	DefaultExchangeType = "fanout"
	DefaultRoutingKey   = "megam_routingkey"
	DefaultConsumerTag  = "megam_node_consumer"
	
)

var (
	mut            sync.Mutex // for conn access
	timeoutRegexp  = regexp.MustCompile(`(TIMED_OUT|timeout)$`)
	notFoundRegexp = regexp.MustCompile(`not found$`)
)

func (b *rabbitmqQ) Get(timeout time.Duration) (*Message, error) {
	return nil, errors.New("Get: Not supported, Handler.start(), subscribe for RabbitMQ.")

}

func (b *rabbitmqQ) Put(m *Message, delay time.Duration) error {
	cons, err := connection()
	if err != nil {
		return err
	}

	//convert Message to "body" bytes
	var body = m.Args[0]

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)
	
	exchange_conf, _ := config.GetString("amqp:exchange")
	if exchange_conf == "" {
		exchange_conf = DefaultExchange
	}
	routingkey_conf, _ := config.GetString("amqp:routingkey")
	if routingkey_conf == "" {
		routingkey_conf = DefaultRoutingKey
	}
	
	if err = cons.channel.Publish(
		exchange_conf, // publish to an exchange
		routingkey_conf, // routing to 0 or more queues
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}
	return err
}

func (b *rabbitmqQ) Delete(m *Message) error {
	return errors.New("Delete: Not supported for RabbitMQ.")

}

func (b *rabbitmqQ) Release(m *Message, delay time.Duration) error {
	return errors.New("Release: Not supported for RabbitMQ.")
}

type rabbitmqFactory struct{}

func (b rabbitmqFactory) Get(name string) (Q, error) {
	return &rabbitmqQ{name: name}, nil
}

func (b rabbitmqFactory) Handler(f func(*Message), name ...string) (Handler, error) {
				log.Printf("RabbitMQ: Handler ")

	return &executor{
		inner: func() {
			if deliveries, err := consume(5e9); err == nil {

				for d := range deliveries {
					log.Printf("got %dB delivery: [%v] %q", len(d.Body), d.DeliveryTag, d.Body)
				    message := &Message{}
					//We have the message here (oo not yet), what do you want to do ?
					//Associate it with a command, and pass it in a go routine ?
					//				log.Printf("Dispatching %q message to handler function.", message.Action)
					go func(m *Message) {
						f(m)
						q := rabbitmqQ{}
						if m.delete {
							q.Delete(m)
						} else {
							q.Release(m, 0)
						}
					}(message)
				}
				log.Printf("handle: deliveries channel closed")
				//done <- nil
			} else {
				log.Printf("Failed to get message from the queue: %s. Trying again...", err)
				if e, ok := err.(*net.OpError); ok && e.Op == "dial" {
					time.Sleep(5e9)
				}

			}
		},
	}, nil
}

func connection() (*Consumer, error) {
	var (
		addr string
		err  error
	)

	mut.Lock()

	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     DefaultConsumerTag,
		done:    make(chan error),
	}

	if c.conn == nil {
		mut.Unlock()
		addr, err = config.GetString("amqp:url")
		if err != nil {
			addr = "localhost:5672"
		}
		mut.Lock()
		if c.conn, err = amqp.Dial(addr); err != nil {
			mut.Unlock()
			return nil, err
		}
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		mut.Unlock()
		return nil, err
	}
	
	exchange_conf, _ := config.GetString("amqp:exchange")
	if exchange_conf == "" {
		exchange_conf = DefaultExchange
	}	

	if err = c.channel.ExchangeDeclare(
		exchange_conf, // name of the exchange
		DefaultExchangeType, // exchange Type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	); err != nil {
		mut.Unlock()
		return nil, err
	}

	mut.Unlock()
	return c, err
}

func rconnection() (*Consumer, error) {
	cons, err := connection()
	mut.Lock()
	
	queue_conf, _ := config.GetString("amqp:queue")
	if queue_conf == "" {
		queue_conf = DefaultQueue
	}
	
	decl_q, err := cons.channel.QueueDeclare(
		queue_conf, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	
	if err != nil {
		mut.Unlock()
		return nil, err
	}
	
	cons.queue = &decl_q
		
	exchange_conf, _ := config.GetString("amqp:exchange")
	if exchange_conf == "" {
		exchange_conf = DefaultExchange
	}
	routingkey_conf, _ := config.GetString("amqp:routingkey")
	if routingkey_conf == "" {
		routingkey_conf = DefaultRoutingKey
	}

	if err = cons.channel.QueueBind(
		cons.queue.Name, // name of the queue
		routingkey_conf,
		exchange_conf,
		false, // noWait
		nil,   // arguments
	); err != nil {
		mut.Unlock()
		return nil, err
	}
	mut.Unlock()
	return cons, nil
}

//returns AMQP Consumer (ASynchronous, blocked - dies on shutdown)
func consume(timeout time.Duration) (<-chan amqp.Delivery, error) {
	cons, err := rconnection()

	deliveries, err := cons.channel.Consume(
		cons.queue.Name, // name
		cons.tag,        // consumerTag,
		false,               // noAck
		false,               // exclusive
		false,               // noLocal
		false,               // noWait
		nil,                 // arguments
	)
	
	if err != nil {
		return nil, err
	}
    
	return deliveries, nil
}

/*
//shut it down, the handler actually shuts it down.
func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

*/
