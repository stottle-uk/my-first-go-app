package linkstatusapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// LinkStatusAPI : LinkStatusAPI
type LinkStatusAPI struct {
	hub *wshub.Hub
}

type options struct {
	Hub *wshub.Hub
}

// LinkStatus : LinkStatus
func LinkStatus(router *mux.Router, hub *wshub.Hub) {
	linkStatusAPI, err := newAPI(options{
		Hub: hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	linkStatusRouter := router.PathPrefix("/link-status").Subrouter()
	linkStatusRouter.HandleFunc("", linkStatusAPI.AddLink).Methods("POST")
}

// AddLink : AddLink
func (s *LinkStatusAPI) AddLink(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	go func() {
		s.hub.SendRestricted <- wshub.Message{ClientIds: []string{"cliendId123"}, Data: body}
	}()

	w.WriteHeader(201)
}

func newAPI(options options) (*LinkStatusAPI, error) {
	s := &LinkStatusAPI{
		hub: options.Hub,
	}
	return s, nil
}
