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
	clients        map[string]*client
	register       chan *client
	unregister     chan *client
	Received       chan []byte
	Broadcast      chan []byte
	SendRestricted chan *Message
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
		client := newClient(id, w, r)

		hub.register <- client

		go client.readPump()
		go client.writePump()
		go hub.handleReceived(client)
	}
}

func newHub() *Hub {
	return &Hub{
		clients:        make(map[string]*client),
		register:       make(chan *client),
		unregister:     make(chan *client),
		Received:       make(chan []byte),
		Broadcast:      make(chan []byte),
		SendRestricted: make(chan *Message),
	}
}

func newClient(id string, w http.ResponseWriter, r *http.Request) *client {
	conn, _ := upgrader.Upgrade(w, r, nil)
	return &client{
		id:       id,
		conn:     conn,
		send:     make(chan []byte),
		received: make(chan []byte)}
}

func (h *Hub) handleReceived(client *client) {
	for {
		select {
		case message := <-client.received:
			h.Received <- message
		}
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

		case message := <-h.SendRestricted:
			for _, id := range message.ClientIds {
				client, ok := h.clients[id]
				if ok {
					select {
					case client.send <- message.Data:
					default:
						delete(h.clients, client.id)
						close(client.send)
					}
				}
			}
		}
	}
}
