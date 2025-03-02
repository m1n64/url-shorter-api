package utils

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
)

const (
	AnalyticsQueue = "analytics"
)

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

type Consumer interface {
	Handle(msg amqp.Delivery)
}

var instance *RabbitMQConnection
var once sync.Once

func ConnectRabbitMQ(rabbitHost string) *RabbitMQConnection {
	once.Do(func() {
		conn, err := amqp.Dial(rabbitHost)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
		}

		instance = &RabbitMQConnection{
			Connection: conn,
			Channel:    ch,
		}
		log.Println("RabbitMQ connection established")
	})
	return instance
}

func InitializeQueues() {
	queues := []string{
		AnalyticsQueue,
	}

	for _, queue := range queues {
		_, err := instance.Channel.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to declare queue %s: %s", queue, err)
		}
		log.Printf("Queue declared: %s", queue)
	}
}

func GetRabbitMQInstance() *RabbitMQConnection {
	if instance == nil {
		log.Fatalf("RabbitMQ connection is not initialized. Call ConnectRabbitMQ first.")
	}
	return instance
}

func (r *RabbitMQConnection) CloseRabbitMQ() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
	log.Println("RabbitMQ connection closed")
}

func ListenToQueue(queueName string, consumer Consumer) error {
	msgs, err := instance.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		go consumer.Handle(msg)
	}

	log.Printf("Started listening to queue: %s", queueName)
	return nil
}
