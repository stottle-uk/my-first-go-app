package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	linkstatus "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus"
	scannertasks "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks"
	websocket "github.com/stottle-uk/my-first-go-app/internal/features/websocket"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

func main() {
	flag.Parse()
	router := mux.NewRouter()
	hub, handler := wshub.CreateHub()

	scannertasks.Init(subRouter(router, "/scanner-tasks"), hub)
	linkstatus.Init(subRouter(router, "/link-status"), hub)
	websocket.Init(subRouter(router, "/ws"), handler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func subRouter(router *mux.Router, path string) *mux.Router {
	return router.PathPrefix(path).Subrouter()
}
