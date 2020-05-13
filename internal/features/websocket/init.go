package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Init : Init
func Init(router *mux.Router, handler func(w http.ResponseWriter, r *http.Request)) {
	router.HandleFunc("/{id}", handler)
}
