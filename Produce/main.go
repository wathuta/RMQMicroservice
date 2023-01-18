package main

import (
	"RMQLogsConsumer/handlers"
	rmqconfig "RMQLogsConsumer/rmqConfig"
	"log"
	"net/http"
)

const (
	URL          = "amqp://guest:guest@localhost:5672/"
	EXCHANGENAME = "log_topic"
	ROUTINGKEY   = "logs.error"
	ROUTINGKEY1  = "logs.info"

	QUEUENAME = "logs_topic"
)
const (
	CONSUMERQUEUE      = "CONSUMERRESPONSEQUEUE"
	RESPONSEROUTINGKEY = "response.logs.after"
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
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	for i := 0; i < 1000; i++ {

// 		if i%2 == 0 {
// 			ch.PublishWithContext(ctx, EXCHANGENAME, ROUTINGKEY, false, false, amqp.Publishing{
// 				ContentType: "text/plain",
// 				Body:        []byte(fmt.Sprintf("body %d", i)),
// 			})
// 			log.Println(i)
// 		} else {
// 			ch.PublishWithContext(ctx, EXCHANGENAME, ROUTINGKEY1, false, false, amqp.Publishing{
// 				ContentType: "text/plain",
// 				Body:        []byte(fmt.Sprintf("body %d", i)),
// 			})
// 			log.Println(i)
// 		}
// 		time.Sleep(3 * time.Second)
// 	}
// 	log.Printf(" [x] Sent %s", "Done")
// }

func main() {
	mux := http.NewServeMux()
	newRmq := rmqconfig.SetUpPublisher()
	newRmq.RabbitMQConsumer()
	mux.Handle("/", handlers.NewUserHandler(newRmq))
	log.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// func main() {
// 	conn, err := amqp.Dial(URL)
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to create a channel")
// 	//Direct binding can be a fanout with a routing key
// 	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "topic", true, false, false, false, nil), "FAILED TO declare exchange")

// 	queue, err := ch.QueueDeclare(CONSUMERQUEUE, true, false, false, false, nil)
// 	failOnError(err, "failed to declare queue")
// 	failOnError(ch.QueueBind(queue.Name, RESPONSEROUTINGKEY, EXCHANGENAME, false, nil), "unable to bind queue")

// 	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
// 	failOnError(err, "Failed to register a consumer")

// 	log.Println(queue.Name, CONSUMERQUEUE)
// 	errChan := make(chan error)

// 	for d := range msgs {
// 		log.Println("Consuming", msgs)
// 		go func(d amqp.Delivery, errChan chan error) {
// 			// errChan <- c.handleConsumer(d)
// 			log.Printf("%s", d.Body)
// 			d.Ack(false)
// 		}(d, errChan)
// 	}

// }
