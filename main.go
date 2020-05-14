package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	linkstatus "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus"
	scannertasks "github.com/stottle-uk/my-first-go-app/internal/features/scannerTasks"
	websocket "github.com/stottle-uk/my-first-go-app/internal/features/websocket"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func main() {
	flag.Parse()
	router := mux.NewRouter()
	hub, handler := wshub.CreateHub()

	scannertasks.New(subRouter(router, "/scanner-tasks"), hub)
	linkstatus.New(subRouter(router, "/link-status"), hub)
	websocket.New(subRouter(router, "/ws"), handler)

	id, err := storageStuff()
	if err != nil {
		log.Printf("Insert Error: %v", err)
	}
	log.Printf("Insert id: %v", id)

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

func storageStuff() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cs, err := connstring.Parse("mongodb://localhost:27017/testing123")
	if err != nil {
		return "", err
	}
	mongoDatabase := cs.Database

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cs.String()))
	defer client.Disconnect(ctx)

	collection := client.Database(mongoDatabase).Collection("linkstatus")

	result, err := collection.InsertOne(ctx, bson.D{
		{"item", "canvas"},
		{"qty", 100},
		{"tags", bson.A{"cotton"}},
		{"size", bson.D{
			{"h", 28},
			{"w", 35.5},
			{"uom", "cm"},
		}},
	})
	if err != nil {
		return "", err
	}

	itemID := result.InsertedID.(primitive.ObjectID).Hex()

	return itemID, nil
}
