package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
)

// New : New
func New(router *mux.Router, handler func(w http.ResponseWriter, r *http.Request)) {
	router.HandleFunc("/{id}", handler)
}
