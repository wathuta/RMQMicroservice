package main

import (
	"RMQConsumer/handlers"
	rmqconfig "RMQConsumer/rmqConfig"
	"RMQConsumer/services"
	"log"
	"net/http"
)

const (
	URL          = "amqp://guest:guest@localhost:5672/"
	EXCHANGENAME = "log_topic"
	ROUTINGKEY   = "logs.error"
	ROUTINGKEY1  = "logs.info"

	QUEUENAME  = "logs_consumer"
	QUEUENAME2 = "test_logs2"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// func main() {
// 	conn, err := amqp.Dial(URL)
// 	failOnError(err, "Failed to connect to RabbitMQ")

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to create a channel")
// 	//Direct binding can be a fanout with a routing key
// 	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "topic", true, false, false, false, nil), "FAILED TO declare exchange")

// 	queue, err := ch.QueueDeclare(QUEUENAME, true, false, false, false, nil)
// 	failOnError(err, "failed to declare queue")
// 	failOnError(ch.QueueBind(queue.Name, ROUTINGKEY, EXCHANGENAME, false, nil), "unable to bind queue")

// 	log.Println(queue.Name)
// 	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
// 	failOnError(err, "Failed to register a consumer")
// 	var forever chan struct{}

// 	for d := range msgs {
// 		go func(d amqp.Delivery) {
// 			log.Printf(" [x] %s", d.Body)
// 			d.Ack(false)
// 		}(d)
// 	}

// 	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
// 	<-forever
// }

func main() {
	newRmq := rmqconfig.NewConsumer()
	services.NewUserService(newRmq)
	go newRmq.Consume()
	go newRmq.SetUpPublisher()

	mux := http.NewServeMux()
	mux.Handle("/", handlers.NewUserHandler())
	log.Println("listening on port :8081")
	http.ListenAndServe(":8081", mux)
}
