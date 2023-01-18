package rmqconfig

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const RESPONSEROUTINGKEY = "response.logs.after"

func (c *ConsumerRabbitMQ) SetUpPublisher() {
	ch, err := c.conn.Channel()
	failOnError(err, "Failed to create a channel")

	failOnError(ch.ExchangeDeclare(EXCHANGENAME, "topic", true, false, false, false, nil), "FAILED TO declare exchange")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c.pubChan = *ch
	c.TimoutCtx = ctx
}
func (r *ConsumerRabbitMQ) Publish(ctx context.Context, byteArr []byte) {
	if err := r.pubChan.PublishWithContext(ctx, EXCHANGENAME, RESPONSEROUTINGKEY, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        byteArr,
	}); err != nil {
		log.Fatal(err)
	}
	log.Println("Published successfully")
}
