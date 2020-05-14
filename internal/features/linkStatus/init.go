package linkstatus

import (
	"fmt"

	"github.com/gorilla/mux"
	api "github.com/stottle-uk/my-first-go-app/internal/features/linkStatus/api"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// New : New
func New(router *mux.Router, hub *wshub.Hub) {
	linkStatusAPI, err := api.NewAPI(api.Options{
		Hub: hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	router.HandleFunc("", linkStatusAPI.AddLink).Methods("POST")
}
