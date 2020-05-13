package scannertask

import (
	"io/ioutil"
	"log"
	"net/http"

	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// API : API
type API struct {
	hub *wshub.Hub
}

// Options : Options
type Options struct {
	Hub *wshub.Hub
}

// NewAPI : NewAPI
func NewAPI(options Options) (*API, error) {
	s := &API{
		hub: options.Hub,
	}

	go s.handleReceived()

	return s, nil
}

// UpdateTask : UpdateTask
func (s *API) UpdateTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	go func() {
		s.hub.Broadcast <- body
	}()

	w.WriteHeader(201)
}

func (s *API) handleReceived() {
	for {
		select {
		case message := <-s.hub.Received:
			s.hub.Broadcast <- message
		}
	}
}
