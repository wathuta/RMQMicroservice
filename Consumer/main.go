package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	URL          = "amqp://guest:guest@localhost:5672/"
	EXCHANGENAME = "New"
	ROUTINGKEY   = "logs"
	ROUTINGKEY1  = "log"

	QUEUENAME  = "test_logs"
	QUEUENAME2 = "test_logs2"
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

	queue, err := ch.QueueDeclare(QUEUENAME, true, false, false, false, nil)
	failOnError(err, "failed to declare queue")
	queue2, err := ch.QueueDeclare(QUEUENAME2, true, false, false, false, nil)
	failOnError(err, "Unable to declare queue")
	failOnError(ch.QueueBind(queue2.Name, ROUTINGKEY1, EXCHANGENAME, false, nil), "unable to bind queue")
	failOnError(ch.QueueBind(queue.Name, ROUTINGKEY, EXCHANGENAME, false, nil), "unable to bind queue")

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	msgs2, err := ch.Consume(QUEUENAME2, "", false, false, false, false, nil)
	failOnError(err, "unable to consumefrom queue")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			d.Ack(false)
		}

	}()
	go func() {
		for d := range msgs2 {
			log.Printf(" [x] %s", d.Body)
			d.Ack(false)
		}

	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
