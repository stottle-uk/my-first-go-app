package scannertasksapi

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

// ScannerTasksAPI : ScannerTasksAPI
type ScannerTasksAPI struct {
	hub *wshub.Hub
}

// NewAPI : NewAPI
func NewAPI(options Options) (*ScannerTasksAPI, error) {
	s := &ScannerTasksAPI{
		hub: options.Hub,
	}
	return s, nil
}

// UpdateTask : UpdateTask
func (s *ScannerTasksAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
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
