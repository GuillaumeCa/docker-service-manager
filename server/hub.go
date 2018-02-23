package main

import "github.com/gorilla/websocket"

// Hub handles websocket connections
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// Client holds the ws conn
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) broadcastJSON(json interface{}) {
	for c := range h.clients {
		c.sendJSON(json)
	}
}

func newClient(hub *Hub, conn *websocket.Conn) *Client {
	client := &Client{hub: hub, conn: conn}
	client.hub.register <- client
	return client
}

func (c *Client) sendJSON(json interface{}) error {
	return c.conn.WriteJSON(json)
}
