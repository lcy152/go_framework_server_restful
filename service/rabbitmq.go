package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"tumor_server/db"

	"github.com/streadway/amqp"
)

type Exchange struct {
	Name string
	Done chan int
}

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	ExchangeMap  map[string]*Exchange
	SendQueue    string
	ReceiveQueue string
	Mqurl        string
	QueueConfig  map[string]interface{}
	Wg           sync.WaitGroup
}

func NewRabbitMQ(sc *Container, url string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		Mqurl:        url,
		SendQueue:    "_tumor_send",
		ReceiveQueue: "_dipper_send",
		sc:           sc,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		log.Println(err)
		return nil
	}
	rabbitmq.ExchangeMap = make(map[string]*Exchange)
	rabbitmq.Connect()
	return rabbitmq
}

func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) Connect() {
	op := db.NewOptions()
	insList, _, err := r.sc.DB.LoadInstitution(context.TODO(), op)
	if err != nil {
		return
	}
	for _, v := range r.ExchangeMap {
		v.Done <- 0
	}
	r.ExchangeMap = make(map[string]*Exchange)
	for _, v := range insList {
		ex := &Exchange{Name: v.Guid}
		ex.Done = make(chan int, 1)
		r.ExchangeMap[v.Guid] = ex
	}
	for _, v := range r.ExchangeMap {
		err := r.channel.ExchangeDeclare(
			v.Name,   // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			fmt.Println(err)
		}
		_, err = r.channel.QueueDeclare(
			v.Name+r.SendQueue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		}
		_, err = r.channel.QueueDeclare(
			v.Name+r.ReceiveQueue,
			true,
			false,
			false,
			false,
			r.QueueConfig,
		)
		if err != nil {
			fmt.Println(err)
		}
		err = r.channel.QueueBind(
			v.Name+r.SendQueue,
			r.SendQueue,
			v.Name,
			false,
			r.QueueConfig,
		)
		if err != nil {
			fmt.Println(err)
		}
		err = r.channel.QueueBind(
			v.Name+r.ReceiveQueue,
			r.ReceiveQueue,
			v.Name,
			false,
			r.QueueConfig,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
	r.Wg.Wait()
	r.Consume()
}

func (r *RabbitMQ) Publish(exchangeName, message string) {
	if _, ok := r.ExchangeMap[exchangeName]; ok {
		log.Println("publish: ", exchangeName)
		err := r.channel.Publish(
			exchangeName,
			r.SendQueue,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Timestamp:    time.Now(),
				Body:         []byte(message),
			},
		)
		if err != nil {
			log.Println(err)
		}
	}
}

func (r *RabbitMQ) Consume() {
	for _, v := range r.ExchangeMap {
		msgs, err := r.channel.Consume(
			v.Name+r.ReceiveQueue,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		}
		messageReceiveFunc := func(v *Exchange) {
			r.Wg.Add(1)
			defer r.Wg.Done()
			for {
				select {
				case d, ok := <-msgs:
					if !ok {
						return
					}
					if HandleMessage(d.Body, d.Exchange) {
						d.Ack(false)
					} else {
						d.Reject(false)
					}
				case <-v.Done:
					log.Printf("message handle done: %s", v.Name)
					return
				}
			}
		}
		go messageReceiveFunc(v)
	}

}
