package main

import (
	"context"
	"log"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 100; i++ {
		err = ch.PublishWithContext(ctx,
			"logs", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("body"),
			})
		failOnError(err, "Failed to publish a message")
		time.Sleep(3 * time.Second)
	}

	log.Printf(" [x] Sent %s", "Done")
}
