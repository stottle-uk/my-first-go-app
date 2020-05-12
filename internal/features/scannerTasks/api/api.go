package scannertasksapi

import (
	"io/ioutil"
	"log"
	"net/http"
)

// ScannerTasksAPI : ScannerTasksAPI
type ScannerTasksAPI struct {
	hub *Hub
}

// NewAPI : NewAPI
func NewAPI() (*ScannerTasksAPI, error) {
	hub := newHub()
	go hub.run()

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
		s.hub.broadcast <- []byte(body)
	}()

	w.WriteHeader(201)
}

// WS : WS
func (s *ScannerTasksAPI) WS(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{id: conn.RemoteAddr().String(), hub: s.hub, conn: conn, send: make(chan []byte)}
	client.hub.register <- client

	go client.writePump()
}
