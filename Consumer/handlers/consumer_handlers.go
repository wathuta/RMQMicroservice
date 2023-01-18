package handlers

import (
	"net/http"
	"regexp"
)

var (
	listUserRe   = regexp.MustCompile(`^\/users[\/]*$`)
	getUserRe    = regexp.MustCompile(`^\/users\/(\d+)$`)
	createUserRe = regexp.MustCompile(`^\/users[\/]*$`)
)

type userHandler struct {
}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet:
		h.LIST(w, r)
	case r.Method == http.MethodPost:
		h.CREATE(w, r)
	case r.Method == http.MethodPut:
		h.UPDATE(w, r)
	case r.Method == http.MethodDelete:
		h.DELETE(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}
func (h *userHandler) CREATE(w http.ResponseWriter, r *http.Request) {
}
func (h *userHandler) LIST(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))

}
func (h *userHandler) UPDATE(w http.ResponseWriter, r *http.Request) {
}
func (h *userHandler) DELETE(w http.ResponseWriter, r *http.Request) {
}
