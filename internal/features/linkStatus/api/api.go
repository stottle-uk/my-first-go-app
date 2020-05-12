package linkstatusapi

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// LinkStatusAPI : LinkStatusAPI
type LinkStatusAPI struct {
	hub *Hub
}

// NewAPI : NewAPI
func NewAPI() (*LinkStatusAPI, error) {
	hub := newHub()
	go hub.run()

	s := &LinkStatusAPI{
		hub: hub,
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
		s.hub.sendRestricted <- Message{ids: []string{"cliendId123"}, data: []byte(body)}
	}()

	w.WriteHeader(201)
}

// WS : WS
func (s *LinkStatusAPI) WS(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	log.Printf("id: %v", id)

	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{id: id, hub: s.hub, conn: conn, send: make(chan []byte)}
	client.hub.register <- client

	go client.writePump()
}
