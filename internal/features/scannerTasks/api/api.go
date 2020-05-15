package scannertask

import (
	"io/ioutil"
	"log"
	"net/http"

	hub "github.com/stottle-uk/my-first-go-app/internal/features/hub"
)

// API : API
type API struct {
	hub *hub.Hub
}

// Options : Options
type Options struct {
	Hub *hub.Hub
}

// NewAPI : NewAPI
func NewAPI(options Options) (*API, error) {
	s := &API{
		hub: options.Hub,
	}

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

	s.hub.Broadcast <- body

	w.WriteHeader(201)
}

// HandleReceived : HandleReceived
func (s *API) HandleReceived() {
	for {
		select {
		case message := <-s.hub.Received:
			s.hub.Broadcast <- message
		}
	}
}
