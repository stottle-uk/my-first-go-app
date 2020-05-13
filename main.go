package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	linkstatus "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	scannertasks "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

func main() {
	flag.Parse()
	router := mux.NewRouter()
	hub, handler := wshub.CreateHub()

	scannertasks.Init(router, hub)
	linkstatus.Init(router, hub)
	webSockets(router, handler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func webSockets(router *mux.Router, handler func(w http.ResponseWriter, r *http.Request)) {
	router.HandleFunc("/ws/{id}", handler)
}
