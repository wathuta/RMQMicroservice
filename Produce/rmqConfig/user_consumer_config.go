package rmqconfig

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	CONSUMERQUEUE      = "CONSUMERRESPONSEQUEUE"
	RESPONSEROUTINGKEY = "response.logs.after"
)

func (c *RabbitMQPublisher) RabbitMQConsumer() {
	ch, err := c.conn.Channel()
	failOnError(err, "Failed to create a channel")
	//Direct binding can be a fanout with a routing key
	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "topic", true, false, false, false, nil), "FAILED TO declare exchange")

	queue, err := ch.QueueDeclare(CONSUMERQUEUE, true, false, false, false, nil)
	failOnError(err, "failed to declare queue")
	failOnError(ch.QueueBind(queue.Name, RESPONSEROUTINGKEY, EXCHANGENAME, false, nil), "unable to bind queue")

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	errChan := make(chan error)

	go func(errChan chan error) {
		for d := range msgs {
			go func(d amqp.Delivery, errChan chan error) {
				// errChan <- c.handleConsumer(d)
				log.Printf("%s", d.Body)
				d.Ack(false)
			}(d, errChan)
		}
	}(errChan)
}
