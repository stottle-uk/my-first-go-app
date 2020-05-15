package scannertasks

import (
	"fmt"

	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
	api "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks/api"
	router "github.com/stottle-uk/my-first-go-app/internal/services/router"
)

// Options : Options
type Options struct {
	Router *router.Router
	Hub    *hub.Hub
}

// New : New
func New(options Options) {
	scannertasksAPI, err := api.NewAPI(api.Options{
		Hub: options.Hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	options.Router.HandlePost("/{id:[0-9]+}", scannertasksAPI.UpdateTask)

	go scannertasksAPI.HandleReceived()
}
