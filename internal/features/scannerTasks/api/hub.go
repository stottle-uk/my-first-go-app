package scannertasksapi

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Message : Message
type Message struct {
	ids  []string
	data []byte
}

// Client : Client
type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub : Hub
type Hub struct {
	clients        map[string]*Client
	broadcast      chan []byte
	sendRestricted chan Message
	register       chan *Client
	unregister     chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newHub() *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		sendRestricted: make(chan Message),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[string]*Client),
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

		case message := <-h.broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client.id)
					close(client.send)
				}
			}

		case m := <-h.sendRestricted:
			for _, id := range m.ids {
				client, ok := h.clients[id]
				if ok {
					select {
					case client.send <- m.data:
					default:
						delete(h.clients, client.id)
						close(client.send)
					}
				}
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
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
