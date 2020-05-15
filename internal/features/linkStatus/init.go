package linkstatus

import (
	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
	api "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	domains "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/domains"
	routes "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/routes"
	store "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/storage"
	redirect "github.com/stottle-uk/my-first-go-app/internal/services/redirect"
	router "github.com/stottle-uk/my-first-go-app/internal/services/router"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"
)

// Hub : Hub
type Hub struct {
	base         *hub.Hub
	tasks        *storage.Collection
	SendByTaskID chan *domains.LinkStatus
}

// Options : Options
type Options struct {
	Router   *router.Router
	Hub      *hub.Hub
	Db       *storage.Database
	Redirect *redirect.Redirect
}

// GetTaskQuery : GetTaskQuery
type GetTaskQuery struct {
	TaskID int
}

// New : New
func New(options Options) {
	store := store.New(options.Db, options.Hub)

	api := api.NewAPI(api.Options{
		Store:    store,
		Redirect: options.Redirect,
	})

	routes := routes.NewRoutes(routes.Options{
		Router: options.Router,
		API:    api,
	})

	routes.UseRoutes()
}
