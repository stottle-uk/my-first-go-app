package websocket

import (
	"net/http"

	router "github.com/stottle-uk/my-first-go-app/internal/services/router"
)

// Options : Options
type Options struct {
	Router  *router.Router
	Handler func(w http.ResponseWriter, r *http.Request)
}

// New : New
func New(options Options) {
	options.Router.Handle("/{id}", options.Handler)
}
