package rmqconfig

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	URL          = "amqp://guest:guest@localhost:5672/"
	EXCHANGENAME = "log_topic"
	ROUTINGKEY   = "logs.error"
	ROUTINGKEY1  = "logs.info"

	QUEUENAME = "logs_topic"
)
const (
	TypeGet      = "GET"
	TypeRegister = "REGISTER"
	TypeUpdate   = "UPDATE"
	TypeDelete   = "DELETE"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type RabbitMQPublisher struct {
	ctx  context.Context
	ch   *amqp.Channel
	conn amqp.Connection
}

func SetUpPublisher() *RabbitMQPublisher {
	conn, err := amqp.Dial(URL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to create a channel")
	//Direct binding can be a fanout with a routing key
	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "topic", true, false, false, false, nil), "FAILED TO declare exchange")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return &RabbitMQPublisher{
		ctx:  ctx,
		ch:   ch,
		conn: *conn,
	}
}
func (r *RabbitMQPublisher) Publish(ctx context.Context, byteArr []byte) {
	if err := r.ch.PublishWithContext(ctx, EXCHANGENAME, ROUTINGKEY, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        byteArr,
		Type:        TypeGet,
	}); err != nil {
		log.Fatal(err)
	}
	time.Sleep(30 * time.Microsecond)
	log.Println("Published successfully")
}
