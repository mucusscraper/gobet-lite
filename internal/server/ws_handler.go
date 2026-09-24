package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/mucusscraper/gobet-lite/internal/websockethub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Liberado para testes locais
}

func ServeWs(hubInstance *websockethub.Hub, w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil || userID <= 0 {
		http.Error(w, "user id invalido para websocket", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	hubInstance.Register(userID, conn)
	go func() {
		defer hubInstance.Unregister(userID)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
