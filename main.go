package main

import (
	"flag"
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

	scannertasksapi.ScannerTasks(router, hub)
	linkstatusapi.LinkStatus(router, hub)
	webSockets(router, handler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func webSockets(router *mux.Router, handler func(w http.ResponseWriter, r *http.Request)) {
	router.HandleFunc("/ws/{id}", handler)
}
