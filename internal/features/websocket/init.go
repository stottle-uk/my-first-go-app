package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Options : Options
type Options struct {
	Router  *mux.Router
	Handler func(w http.ResponseWriter, r *http.Request)
}

// New : New
func New(options Options) {
	options.Router.HandleFunc("/{id}", options.Handler)
}
