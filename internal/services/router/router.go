package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Router : Hub
type Router struct {
	base *mux.Router
}

// New : New
func New() *Router {
	return &Router{
		base: mux.NewRouter(),
	}
}

// SubRouter : SubRouter
func (r *Router) SubRouter(path string) *Router {
	return &Router{
		base: r.base.PathPrefix(path).Subrouter(),
	}
}

// UseCors : UseCors
func (r *Router) UseCors() http.Handler {
	cr := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"},
		AllowCredentials: true,
		Debug:            false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	})
	return cr.Handler(r.base)
}

// Handle : Handle
func (r *Router) Handle(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.handle(path, handler)
}

// HandlePost : HandlePost
func (r *Router) HandlePost(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.handle(path, handler).Methods("POST")
}

func (r *Router) handle(path string, handler func(w http.ResponseWriter, r *http.Request)) *mux.Route {
	return r.base.HandleFunc(path, handler)
}
