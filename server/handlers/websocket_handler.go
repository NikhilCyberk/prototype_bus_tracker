// handlers/websocket_handler.go
package handlers

import (
	"bustracking/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

type WebSocketHandler struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan models.Message
	Mutex     sync.Mutex
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan models.Message),
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	// Register client
	h.Mutex.Lock()
	h.Clients[ws] = true
	h.Mutex.Unlock()

	// Read messages from the client
	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			h.Mutex.Lock()
			delete(h.Clients, ws)
			h.Mutex.Unlock()
			break
		}

		// Handle the message based on its type
		switch msg.Type {
		case "subscribe":
			// Handle subscription requests
			log.Println("Client subscribed to updates")
		case "location_update":
			// Process location update from drivers
			h.Broadcast <- msg
		}
	}
}

func (h *WebSocketHandler) BroadcastMessages() {
	for {
		msg := <-h.Broadcast

		h.Mutex.Lock()
		for client := range h.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error sending message:", err)
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mutex.Unlock()
	}
}
