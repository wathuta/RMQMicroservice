package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	URL          = "amqp://guest:guest@localhost:5672/"
	EXCHANGENAME = "New"
	ROUTINGKEY   = "logs"
	ROUTINGKEY1  = "log"

	QUEUENAME = "test_logs"
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
	//Direct binding can be a fanout with a routing key
	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "direct", true, false, false, false, nil), "FAILED TO declare exchange")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < 1000; i++ {

		if i%2 == 0 {
			ch.PublishWithContext(ctx, EXCHANGENAME, ROUTINGKEY, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("body %d", i)),
			})
			log.Println(i)
		} else {
			ch.PublishWithContext(ctx, EXCHANGENAME, ROUTINGKEY1, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("body %d", i)),
			})
			log.Println(i)
		}
		time.Sleep(3 * time.Second)
	}
	log.Printf(" [x] Sent %s", "Done")
}
