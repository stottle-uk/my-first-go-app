package wshub

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Message : Message
type Message struct {
	ClientIds []string
	Data      []byte
}

// Hub : Hub
type Hub struct {
	clients        map[string]*Client
	register       chan *Client
	unregister     chan *Client
	Broadcast      chan []byte
	SendRestricted chan Message
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CreateHub : CreateHub
func CreateHub() (*Hub, func(w http.ResponseWriter, r *http.Request)) {
	hub := newHub()
	go hub.run()

	return hub, func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		conn, _ := upgrader.Upgrade(w, r, nil)
		client := &Client{id: id, conn: conn, send: make(chan []byte)}
		hub.register <- client

		go client.writePump()
	}
}

func newHub() *Hub {
	return &Hub{
		clients:        make(map[string]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		Broadcast:      make(chan []byte),
		SendRestricted: make(chan Message),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client

		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}

		case message := <-h.Broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client.id)
					close(client.send)
				}
			}

		case m := <-h.SendRestricted:
			for _, id := range m.ClientIds {
				client, ok := h.clients[id]
				if ok {
					select {
					case client.send <- m.Data:
					default:
						delete(h.clients, client.id)
						close(client.send)
					}
				}
			}
		}
	}
}
