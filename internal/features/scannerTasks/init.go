package scannertasks

import (
	"fmt"

	"github.com/gorilla/mux"
	api "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks/api"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// Init : Init
func Init(router *mux.Router, hub *wshub.Hub) {
	scannertasksAPI, err := api.NewAPI(api.Options{
		Hub: hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	scannerTasksRouter := router.PathPrefix("/scanner-tasks").Subrouter()
	scannerTasksRouter.HandleFunc("/{id:[0-9]+}", scannertasksAPI.UpdateTask).Methods("POST")
}
