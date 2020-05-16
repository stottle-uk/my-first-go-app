package linkstatus

import (
	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
	api "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	lsHub "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/hub"
	routes "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/routes"
	store "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/storage"
	redirect "github.com/stottle-uk/my-first-go-app/internal/services/redirect"
	router "github.com/stottle-uk/my-first-go-app/internal/services/router"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"
)

// Options : Options
type Options struct {
	Router   *router.Router
	Hub      *hub.Hub
	Col      *storage.Collection
	Redirect *redirect.Redirect
}

// New : New
func New(options Options) {
	store := store.New(store.Options{
		Links: options.Col,
	})

	hub := lsHub.New(lsHub.Options{
		Hub:   options.Hub,
		Store: store,
	})

	api := api.NewAPI(api.Options{
		Store:    hub,
		Redirect: options.Redirect,
	})

	routes := routes.NewRoutes(routes.Options{
		Router: options.Router,
		API:    api,
	})

	routes.UseRoutes()
}
