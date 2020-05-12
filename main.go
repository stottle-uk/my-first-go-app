package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	linkstatusapi "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	scannertasksapi "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks/api"
)

func main() {
	flag.Parse()
	router := mux.NewRouter()

	scannerTasksAPI(router)
	linkStatusAPI(router)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func scannerTasksAPI(router *mux.Router) {
	scannertasksAPI, err := scannertasksapi.NewAPI()
	if err != nil {
		fmt.Println(err)
	}
	scannerTasksRouter := router.PathPrefix("/scanner-tasks").Subrouter()
	scannerTasksRouter.HandleFunc("/{id:[0-9]+}", scannertasksAPI.UpdateTask).Methods("POST")
	scannerTasksRouter.HandleFunc("/ws", scannertasksAPI.WS)
}

func linkStatusAPI(router *mux.Router) {
	linkStatusAPI, err := linkstatusapi.NewAPI()
	if err != nil {
		fmt.Println(err)
	}
	linkStatusRouter := router.PathPrefix("/link-status").Subrouter()
	linkStatusRouter.HandleFunc("", linkStatusAPI.AddLink).Methods("POST")
	linkStatusRouter.HandleFunc("/ws/{id}", linkStatusAPI.WS)
}
