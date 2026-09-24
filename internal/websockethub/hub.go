package websockethub

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// esse codigo vai ter as funcoes que vao permitir o registro de quem esta online e permite enviar msgs direcionadas
// por exemplo, atualizar o saldo de um usuario especifico

type Client struct {
	UserID int
	Conn   *websocket.Conn
}

type Hub struct {
	clients map[int]*Client
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int]*Client),
	}
}

func (h *Hub) Register(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = &Client{UserID: userID, Conn: conn}
	log.Printf("[Websocket] usuario %d conectado!\n", userID)
}

func (h *Hub) Unregister(userID int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client, ok := h.clients[userID]; ok {
		client.Conn.Close()
		delete(h.clients, userID)
		log.Printf("[Websocket] usuario %d desconectado!\n", userID)
	}
}

func (h *Hub) SendToUser(userID int, payload interface{}) {
	h.mu.Lock()
	client, ok := h.clients[userID]
	h.mu.Unlock()

	if !ok {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("[WebSocket] Erro ao enviar mensagem para usuário %d: %v", userID, err)
		h.Unregister(userID)
	}
}
