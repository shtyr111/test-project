package websocket

import (
	"encoding/json"
	"sync"
	"test-project/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

type WebsocketService struct {
	clients map[uuid.UUID]*websocket.Conn
	mu      sync.Mutex
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{clients: make(map[uuid.UUID]*websocket.Conn)}
}

func (ws *WebsocketService) AddClient(id uuid.UUID, conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.clients[id] = conn
}

func (ws *WebsocketService) RemoveClient(id uuid.UUID) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	delete(ws.clients, id)
}

func (ws *WebsocketService) SendToClient(notification *models.Notification) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	message, err := json.Marshal(notification)
	if err != nil {
		log.Println("Error marshalling price update:", err)
		return
	}

	client, ok := ws.clients[notification.Id]

	if !ok {
		return
	}

	if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println("Error writing message to client:", err)
		client.Close()
		delete(ws.clients, notification.Id)
	}
}
