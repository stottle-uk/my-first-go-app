package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	linkstatusapi "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	scannertasksapi "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks/api"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

func main() {
	flag.Parse()
	router := mux.NewRouter()
	hub, handler := wshub.CreateHub()

	scannerTasksAPI(router, hub)
	linkStatusAPI(router, hub)
	webSockets(router, handler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func scannerTasksAPI(router *mux.Router, hub *wshub.Hub) {
	scannertasksAPI, err := scannertasksapi.NewAPI(scannertasksapi.Options{
		Hub: hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	scannerTasksRouter := router.PathPrefix("/scanner-tasks").Subrouter()
	scannerTasksRouter.HandleFunc("/{id:[0-9]+}", scannertasksAPI.UpdateTask).Methods("POST")
}

func linkStatusAPI(router *mux.Router, hub *wshub.Hub) {
	linkStatusAPI, err := linkstatusapi.NewAPI(linkstatusapi.Options{
		Hub: hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	linkStatusRouter := router.PathPrefix("/link-status").Subrouter()
	linkStatusRouter.HandleFunc("", linkStatusAPI.AddLink).Methods("POST")
}

func webSockets(router *mux.Router, handler func(w http.ResponseWriter, r *http.Request)) {
	router.HandleFunc("/ws/{id}", handler)
}
