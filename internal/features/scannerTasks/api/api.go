package scannertasksapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// ScannerTasksAPI : ScannerTasksAPI
type ScannerTasksAPI struct {
	hub *wshub.Hub
}

type options struct {
	Hub *wshub.Hub
}

// ScannerTasks : ScannerTasks
func ScannerTasks(router *mux.Router, hub *wshub.Hub) {
	scannertasksAPI, err := newAPI(options{
		Hub: hub,
	})
	if err != nil {
		fmt.Println(err)
	}
	scannerTasksRouter := router.PathPrefix("/scanner-tasks").Subrouter()
	scannerTasksRouter.HandleFunc("/{id:[0-9]+}", scannertasksAPI.UpdateTask).Methods("POST")
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

func (s *ScannerTasksAPI) handleReceived() {
	for {
		select {
		case message := <-s.hub.Received:
			s.hub.Broadcast <- message
		}
	}
}

func newAPI(options options) (*ScannerTasksAPI, error) {
	s := &ScannerTasksAPI{
		hub: options.Hub,
	}

	go s.handleReceived()

	return s, nil
}
