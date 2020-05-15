package main

import (
	"flag"
	"net/http"

	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
	linkstatus "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus"
	scannertasks "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks"
	websocket "github.com/stottle-uk/my-first-go-app/internal/features/websocket"
	redirect "github.com/stottle-uk/my-first-go-app/internal/services/redirect"
	router "github.com/stottle-uk/my-first-go-app/internal/services/router"
	storage "github.com/stottle-uk/my-first-go-app/internal/services/storage"
)

func main() {
	flag.Parse()
	router := router.New()
	db := storage.NewDb()
	hub, handler := hub.NewHub(db.Collection("tasks"))
	redirect := redirect.New()

	scannertasks.New(scannertasks.Options{
		Router: router.SubRouter("/scanner-tasks"),
		Hub:    hub})

	linkstatus.New(linkstatus.Options{
		Router:   router.SubRouter("/link-status"),
		Hub:      hub,
		Db:       db,
		Redirect: redirect})

	websocket.New(websocket.Options{
		Router:  router.SubRouter("/ws"),
		Handler: handler})

	http.Handle("/", router.UseCors())
	http.ListenAndServe(":8080", nil)
}
