package scannertasks

import (
	"fmt"

	"github.com/gorilla/mux"
	api "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks/api"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// Options : Options
type Options struct {
	Router *mux.Router
	Hub    *wshub.Hub
}

// New : New
func New(options Options) {
	scannertasksAPI, err := api.NewAPI(api.Options{
		Hub: options.Hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	options.Router.HandleFunc("/{id:[0-9]+}", scannertasksAPI.UpdateTask).Methods("POST")

	go scannertasksAPI.HandleReceived()
}
