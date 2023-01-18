package handlers

import (
	"RMQLogsConsumer/models"
	rmqconfig "RMQLogsConsumer/rmqConfig"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

var (
	listUserRe   = regexp.MustCompile(`^\/users[\/]*$`)
	getUserRe    = regexp.MustCompile(`^\/users\/(\d+)$`)
	createUserRe = regexp.MustCompile(`^\/users[\/]*$`)
)

const (
	TypeGet      = "GET"
	TypeRegister = "REGISTER"
	TypeUpdate   = "UPDATE"
	TypeDelete   = "DELETE"
)

type userHandler struct {
	rmq *rmqconfig.RabbitMQPublisher
}

func NewUserHandler(injectedrmq *rmqconfig.RabbitMQPublisher) *userHandler {
	return &userHandler{rmq: injectedrmq}
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet:
		h.LIST(w, r)
		return
	case r.Method == http.MethodPost:
		h.CREATE(w, r)
		return
	case r.Method == http.MethodPut:
		h.UPDATE(w, r)
		return
	case r.Method == http.MethodDelete:
		h.DELETE(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
		return
	}
}
func (h *userHandler) CREATE(w http.ResponseWriter, r *http.Request) {
}
func (h *userHandler) LIST(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 1000; i++ {
		message := models.Message{
			Type:  TypeGet,
			Value: fmt.Sprintf("hello world %d", i),
		}
		pubMessage, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
		}
		h.rmq.Publish(r.Context(), pubMessage)
		time.Sleep(3 * time.Millisecond)
	}
}
func (h *userHandler) UPDATE(w http.ResponseWriter, r *http.Request) {
}
func (h *userHandler) DELETE(w http.ResponseWriter, r *http.Request) {
}
