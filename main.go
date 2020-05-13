package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	cr := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3001"},
		AllowCredentials: true,
		Debug:            false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	})
	crHandler := cr.Handler(router)
	http.Handle("/", crHandler)
	http.ListenAndServe(":8080", nil)
}

func subRouter(router *mux.Router, path string) *mux.Router {
	return router.PathPrefix(path).Subrouter()
}
