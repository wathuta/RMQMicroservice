package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	URL          = "amqp://guest:guest@localhost:5672/"
	EXCHANGENAME = "logs"
	QUEUENAME    = "logs_consumer"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func main() {
	conn, err := amqp.Dial(URL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to create a channel")

	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "fanout", false, true, false, false, nil), "FAILED TO declare exchange")

	queue, err := ch.QueueDeclare(QUEUENAME, false, true, false, false, nil)
	failOnError(err, "failed to declare queue")

	failOnError(ch.QueueBind(queue.Name, "", EXCHANGENAME, false, nil), "unable to bind queue")
	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			d.Ack(false)
		}

	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
