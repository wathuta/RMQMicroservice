package rmqconfig

import (
	"RMQConsumer/models"
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	URL          = "amqp://guest:guest@localhost:5672/"
	EXCHANGENAME = "log_topic"
	ROUTINGKEY   = "logs.error"
	ROUTINGKEY1  = "logs.info"

	QUEUENAME  = "logs_consumer"
	QUEUENAME2 = "test_logs2"
)

const (
	TypeGet      = "GET"
	TypeRegister = "REGISTER"
	TypeUpdate   = "UPDATE"
	TypeDelete   = "DELETE"
)

type ConsumerRabbitMQ struct {
	conn      amqp.Connection
	pubChan   amqp.Channel
	TimoutCtx context.Context
	Amqpchan  chan models.Message
	Done      chan struct{}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func NewConsumer() *ConsumerRabbitMQ {
	conn, err := amqp.Dial(URL)
	failOnError(err, "Failed to connect to RabbitMQ")

	return &ConsumerRabbitMQ{
		conn:     *conn,
		Amqpchan: make(chan models.Message),
	}
}
func (c *ConsumerRabbitMQ) Consume() {

	ch, err := c.conn.Channel()
	failOnError(err, "Failed to create a channel")
	//Direct binding can be a fanout with a routing key
	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "topic", true, false, false, false, nil), "FAILED TO declare exchange")

	queue, err := ch.QueueDeclare(QUEUENAME, true, false, false, false, nil)
	failOnError(err, "failed to declare queue")
	failOnError(ch.QueueBind(queue.Name, ROUTINGKEY, EXCHANGENAME, false, nil), "unable to bind queue")

	log.Println(queue.Name)
	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	errChan := make(chan error)
	go func(errChan chan error) {
		for d := range msgs {

			err := c.handleConsumer(d)
			if err != nil {
				log.Printf("message  unmarshal error: %s", err.Error())
			}
			d.Ack(false)
		}
	}(errChan)
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	err = <-errChan
	log.Fatal(err, "here")
}

func (c *ConsumerRabbitMQ) handleConsumer(d amqp.Delivery) error {
	received := models.Message{}
	if err := json.Unmarshal(d.Body, &received); err != nil {
		log.Println(err, "here")
		return err
	}

	log.Printf("%s", d.Type)
	switch d.Type {
	case TypeGet:
		c.Amqpchan <- received
	case TypeRegister:
	case TypeDelete:
	case TypeUpdate:

	}

	return nil
}
