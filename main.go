package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	linkstatus "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus"
	scannertasks "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks"
	websocket "github.com/stottle-uk/my-first-go-app/internal/features/websocket"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
	storagemongo "github.com/stottle-uk/my-first-go-app/internal/services/storage"
)

func main() {
	flag.Parse()
	router := mux.NewRouter()
	hub, handler := wshub.CreateHub()
	db, err := storagemongo.NewDb()
	if err != nil {
		log.Printf("Insert Error: %v", err)
	}

	scannertasks.New(subRouter(router, "/scanner-tasks"), hub)
	linkstatus.New(subRouter(router, "/link-status"), hub, db)
	websocket.New(subRouter(router, "/ws"), handler)

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
