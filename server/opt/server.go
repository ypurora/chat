package opt

import (
	"chat/model"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	// Register new client
	clientsMux.Lock()
	clients[ws] = "" // Will be set once we receive the first message
	clientsMux.Unlock()
	log.Printf("Client connected: %v", ws.RemoteAddr())

	// create default room
	createMsg := model.Message{
		Sender:  "server",
		Content: "default",
		Type:    "create_room",
	}
	handleCreateRoom(ws, createMsg)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client %v: %v", ws.RemoteAddr(), err)
			clientsMux.Lock()
			delete(clients, ws)
			clientsMux.Unlock()
			leaveAllRoomsExceptDefault(ws)
			log.Printf("Client disconnected: %v", ws.RemoteAddr())
			break
		}

		//log.Printf("Server received raw message: %s\n", message) // Log the raw message

		var msg model.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error unmarshaling JSON from client %v: %v", ws.RemoteAddr(), err)
			continue
		}

		//log.Printf("Server received message: %+v\n", msg) // Log the parsed message

		// Handle different message types
		switch msg.Type {
		case "auth":
			handleAuth(ws, msg)
		case "create_room":
			handleCreateRoom(ws, msg)
		case "join_room":
			handleJoinRoom(ws, msg)
		case "leave_room":
			handleLeaveRoom(ws, msg)
		default:
			log.Printf("Broadcasting message: %+v\n", msg) // Log the broadcast message
			broadcast <- msg
			storeMessage(msg) // Store message in file
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		if len(msg.Receiver) > 0 && msg.Type != "dm" {
			// Handle room messages
			for _, room := range msg.Receiver {
				msg.Receiver = []string{room}
				handleRoomMessage(msg)
			}
		} else if msg.Type == "dm" {
			handleDirectMessage(msg)
		} else {
			// Broadcast message to all clients
			handleBroadcastMessage(msg)
		}
	}
}

func storeMessage(msg model.Message) {
	if msg.Type == "dm" {
		// Handle direct message storage
		if len(msg.Receiver) > 0 {
			sender := msg.Sender
			receiver := msg.Receiver[0]
			key := generateDMKey(sender, receiver)

			dmDir := filepath.Join(OutputDir, "dm", key)
			if err := os.MkdirAll(dmDir, os.ModePerm); err != nil {
				log.Printf("Error creating directory for DM %s: %v", key, err)
				return
			}

			filePath := filepath.Join(dmDir, "messages.json")
			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Error opening file for DM %s: %v", key, err)
				return
			}
			defer file.Close()

			messageData, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				return
			}
			if _, err := file.Write(append(messageData, '\n')); err != nil {
				log.Printf("Error writing message to file: %v", err)
				return
			}
			log.Printf("Stored DM from %s to %s", sender, receiver)
		}
	} else {
		// Handle room message storage
		for _, room := range msg.Receiver {
			group := room
			if group == "" {
				group = "default"
			}
			username := msg.Sender

			groupDir := filepath.Join(OutputDir, group)
			if err := os.MkdirAll(groupDir, os.ModePerm); err != nil {
				log.Printf("Error creating directory for group %s: %v", group, err)
				return
			}

			filePath := filepath.Join(groupDir, username+".json")
			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Error opening file for user %s in group %s: %v", username, group, err)
				return
			}
			defer file.Close()

			messageData, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				return
			}
			if _, err := file.Write(append(messageData, '\n')); err != nil {
				log.Printf("Error writing message to file: %v", err)
				return
			}
			log.Printf("Stored message from %s in group %s", username, group)
		}
	}
}

func generateDMKey(sender, receiver string) string {
	participants := []string{sender, receiver}
	sort.Strings(participants)
	return strings.Join(participants, "_")
}
