package opt

import (
	"chat/model"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	clients   = make(map[*websocket.Conn]string) // Connected clients and their usernames
	broadcast = make(chan model.Message)         // Broadcast channel
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins, customize for CORS
		},
	}
	clientsMux sync.Mutex // Mutex for concurrent access to clients
	UserDB     = "users.json"
	rooms      = make(map[string]map[*websocket.Conn]bool) // Chat rooms and their members
	roomsMux   sync.Mutex                                  // Mutex for concurrent access to rooms
	OutputDir  = "output"
)
