package services

import (
	rmqconfig "RMQConsumer/rmqConfig"
	"context"
	"log"
)

const (
	TypeGet      = "GET"
	TypeRegister = "REGISTER"
	TypeUpdate   = "UPDATE"
	TypeDelete   = "DELETE"
)

type userService struct {
	rmq *rmqconfig.ConsumerRabbitMQ
}
type UserService interface {
	CREATE(ctx context.Context, req interface{})
	LIST(ctx context.Context, req interface{})
	UPDATE(ctx context.Context, req interface{})
	DELETE(ctx context.Context, req interface{})
}

func NewUserService(injectedRMQ *rmqconfig.ConsumerRabbitMQ) UserService {
	ret := &userService{rmq: injectedRMQ}
	go ret.ROUTE()
	return ret
}

func (u *userService) ROUTE() error {
	for {
		select {
		case received := <-u.rmq.Amqpchan:
			switch {
			case received.Type == TypeGet:
				u.LIST(context.Background(), received.Value)
			case received.Type == TypeDelete:
				u.DELETE(context.Background(), received.Value)
			case received.Type == TypeRegister:
				u.CREATE(context.Background(), received.Value)
			case received.Type == TypeUpdate:
				u.UPDATE(context.Background(), received.Value)
			}
		case <-u.rmq.Done:
			return nil
		}
	}
	// return nil
}

func (u *userService) CREATE(ctx context.Context, req interface{}) {

}
func (u *userService) LIST(ctx context.Context, req interface{}) {
	// for i := 0; i < 1000; i++ {
	// 	message := models.Message{
	// 		Type:  TypeGet,
	// 		Value: fmt.Sprintf("hello world %d", i),
	// 	}
	// 	pubMessage, err := json.Marshal(message)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	u.rmq.Publish(ctx, pubMessage)
	// 	log.Println("Publishing", message)
	// 	time.Sleep(3 * time.Second)
	// }
	log.Println(req, "here")
}
func (u *userService) UPDATE(ctx context.Context, req interface{}) {

}
func (u *userService) DELETE(ctx context.Context, req interface{}) {

}
