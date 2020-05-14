package linkstatus

import (
	"fmt"

	"github.com/gorilla/mux"
	api "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	store "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/storage"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
	"go.mongodb.org/mongo-driver/mongo"
)

// New : New
func New(router *mux.Router, hub *wshub.Hub, db *mongo.Database) {
	store := store.New(db)

	linkStatusAPI, err := api.NewAPI(api.Options{
		Hub:   hub,
		Store: store,
	})
	if err != nil {
		fmt.Println(err)
	}
	router.HandleFunc("", linkStatusAPI.AddLink).Methods("POST")
}
