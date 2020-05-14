package linkstatus

import (
	"fmt"

	"github.com/gorilla/mux"
	api "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	store "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/storage"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
	redirect "github.com/stottle-uk/my-first-go-app/internal/services/redirect"
	"go.mongodb.org/mongo-driver/mongo"
)

// Options : Options
type Options struct {
	Router   *mux.Router
	Hub      *wshub.Hub
	Db       *mongo.Database
	Redirect *redirect.Redirect
}

// New : New
func New(options Options) {
	linkStatusAPI, err := api.NewAPI(api.Options{
		Hub:      options.Hub,
		Store:    store.New(options.Db),
		Redirect: options.Redirect,
	})
	if err != nil {
		fmt.Println(err)
	}
	options.Router.HandleFunc("", linkStatusAPI.AddLink).Methods("POST")
}
