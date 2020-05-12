package scannertasksapi

import (
	"io/ioutil"
	"log"
	"net/http"

	wshub "github.com/stottle-uk/my-first-go-app/internal/services/hub"
)

// ScannerTasksAPI : ScannerTasksAPI
type ScannerTasksAPI struct {
	hub *wshub.Hub
}

// NewAPI : NewAPI
func NewAPI() (*ScannerTasksAPI, error) {
	hub := wshub.NewHub()
	go hub.Run()

	s := &ScannerTasksAPI{
		hub: hub,
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
		s.hub.Broadcast <- []byte(body)
	}()

	w.WriteHeader(201)
}

// WS : WS
func (s *ScannerTasksAPI) WS(w http.ResponseWriter, r *http.Request) {
	conn, _ := wshub.Upgrader.Upgrade(w, r, nil)

	client := &wshub.Client{ID: conn.RemoteAddr().String(), Hub: s.hub, Conn: conn, Send: make(chan []byte)}
	client.Hub.Register <- client

	go client.WritePump()
}
