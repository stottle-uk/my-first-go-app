package linkstatusapi

import (
	"io/ioutil"
	"log"
	"net/http"

	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// Options : Options
type Options struct {
	Hub *wshub.Hub
}

// LinkStatusAPI : LinkStatusAPI
type LinkStatusAPI struct {
	hub *wshub.Hub
}

// NewAPI : NewAPI
func NewAPI(options Options) (*LinkStatusAPI, error) {
	s := &LinkStatusAPI{
		hub: options.Hub,
	}
	return s, nil
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
