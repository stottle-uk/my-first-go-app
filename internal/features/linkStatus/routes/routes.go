package linkstatusapi

import (
	api "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	router "github.com/stottle-uk/my-first-go-app/internal/services/router"
)

// Routes : Routes
type Routes struct {
	router *router.Router
	api    *api.API
}

// Options : Options
type Options struct {
	Router *router.Router
	API    *api.API
}

// NewRoutes : NewRoutes
func NewRoutes(options Options) *Routes {
	return &Routes{router: options.Router, api: options.API}
}

// UseRoutes : UseRoutes
func (r *Routes) UseRoutes() {
	r.router.HandlePost("", r.api.AddLink)
}
