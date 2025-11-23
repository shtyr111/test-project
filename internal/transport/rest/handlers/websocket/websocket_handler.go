package websocket

import (
	"net/http"
	ws "test-project/internal/service/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

type WebsocketHandler struct {
	wsService *ws.WebsocketService
}

func NewWebsocketHandler(wsService *ws.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{wsService: wsService}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (ws WebsocketHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getWebSocketHandler(ws, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getWebSocketHandler(ws WebsocketHandler, w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	id, err := uuid.Parse(strId)
	if err != nil {
		log.Error("Invalid Websocket ID")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Upgrade error:", err)
		return
	}
	defer conn.Close()

	ws.wsService.AddClient(id, conn)
	log.Println("Client connected with ID: ", id)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Client disconnected: %v", err)
			ws.wsService.RemoveClient(id)
			break
		}
	}
}
