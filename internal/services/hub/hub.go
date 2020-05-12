package wshub

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Message : Message
type Message struct {
	ClientIds []string
	Data      []byte
}

// Client : Client
type Client struct {
	ID   string
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

// Hub : Hub
type Hub struct {
	clients        map[string]*Client
	Broadcast      chan []byte
	SendRestricted chan Message
	Register       chan *Client
	Unregister     chan *Client
}

// Upgrader : Upgrader
var Upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewHub : NewHub
func NewHub() *Hub {
	return &Hub{
		Broadcast:      make(chan []byte),
		SendRestricted: make(chan Message),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		clients:        make(map[string]*Client),
	}
}

// Run : Run
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client.ID] = client

		case client := <-h.Unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}

		case message := <-h.Broadcast:
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					delete(h.clients, client.ID)
					close(client.Send)
				}
			}

		case m := <-h.SendRestricted:

			for _, id := range m.ClientIds {
				client, ok := h.clients[id]
				if ok {
					select {
					case client.Send <- m.Data:
					default:
						delete(h.clients, client.ID)
						close(client.Send)
					}
				}
			}
		}
	}
}

// WritePump : WritePump
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
