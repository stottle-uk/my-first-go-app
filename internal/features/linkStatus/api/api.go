package linkstatusapi

import (
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

// NewAPI : NewAPI
func NewAPI() (*LinkStatusAPI, error) {
	hub := wshub.NewHub()
	go hub.Run()

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
		s.hub.SendRestricted <- wshub.Message{ClientIds: []string{"cliendId123"}, Data: []byte(body)}
	}()

	w.WriteHeader(201)
}

// WS : WS
func (s *LinkStatusAPI) WS(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	conn, _ := wshub.Upgrader.Upgrade(w, r, nil)
	client := &wshub.Client{ID: id, Hub: s.hub, Conn: conn, Send: make(chan []byte)}
	client.Hub.Register <- client

	go client.WritePump()
}
