package opt

import (
	"chat/model"
	"github.com/gorilla/websocket"
	"log"
)

// Updated function to leave all rooms except the default room
func leaveAllRoomsExceptDefault(ws *websocket.Conn) {
	roomsMux.Lock()
	defer roomsMux.Unlock()
	for room := range rooms {
		if room != "default" {
			delete(rooms[room], ws)
			if len(rooms[room]) == 0 {
				delete(rooms, room)
			}
		}
	}
	log.Printf("Client %v left all rooms except default", ws.RemoteAddr())
}

func handleCreateRoom(ws *websocket.Conn, msg model.Message) {
	roomsMux.Lock()
	defer roomsMux.Unlock()
	if _, exists := rooms[msg.Content]; exists {
		response := model.Message{
			Sender:  "server",
			Content: "Room already exists",
			Type:    "error",
		}
		ws.WriteJSON(response)
		return
	}
	rooms[msg.Content] = make(map[*websocket.Conn]bool)
	response := model.Message{
		Sender:  "server",
		Content: "Room created",
		Type:    "info",
	}
	ws.WriteJSON(response)
	log.Printf("Room created: %s", msg.Content)
}

func handleJoinRoom(ws *websocket.Conn, msg model.Message) {
	roomsMux.Lock()
	defer roomsMux.Unlock()
	if _, exists := rooms[msg.Content]; !exists {
		response := model.Message{
			Sender:  "server",
			Content: "Room does not exist",
			Type:    "error",
		}
		ws.WriteJSON(response)
		return
	}

	rooms[msg.Content][ws] = true
	response := model.Message{
		Sender:  "server",
		Content: "Joined room",
		Type:    "info",
	}
	ws.WriteJSON(response)
	log.Printf("Client %v joined room: %s", ws.RemoteAddr(), msg.Content)
}

func handleLeaveRoom(ws *websocket.Conn, msg model.Message) {
	roomsMux.Lock()
	defer roomsMux.Unlock()
	if _, exists := rooms[msg.Content]; !exists {
		response := model.Message{
			Sender:  "server",
			Content: "Room does not exist",
			Type:    "error",
		}
		ws.WriteJSON(response)
		return
	}
	delete(rooms[msg.Content], ws)
	if len(rooms[msg.Content]) == 0 {
		delete(rooms, msg.Content)
	}
	response := model.Message{
		Sender:  "server",
		Content: "Left room",
		Type:    "info",
	}
	ws.WriteJSON(response)
	log.Printf("Client %v left room: %s", ws.RemoteAddr(), msg.Content)
}

func handleRoomMessage(msg model.Message) {
	roomsMux.Lock()
	defer roomsMux.Unlock()
	room, exists := rooms[msg.Receiver[0]] // Since we are handling one room at a time in HandleMessages
	if !exists {
		log.Printf("Room %s does not exist", msg.Receiver)
		return
	}
	log.Printf("Broadcasting message to room %s: %+v", msg.Receiver, msg)
	for client := range room {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Error writing JSON to client %v: %v", client.RemoteAddr(), err)
			client.Close()
			delete(room, client)
		}
	}
}

func handleBroadcastMessage(msg model.Message) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	log.Printf("Broadcasting message to all clients: %+v", msg)
	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("Error writing JSON to client %v: %v", client.RemoteAddr(), err)
			client.Close()
			delete(clients, client)
		}
	}
}

func handleDirectMessage(msg model.Message) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	log.Printf("Sending direct message to %s: %+v", msg.Receiver, msg)
	for client, username := range clients {
		if username == msg.Receiver[0] {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error writing JSON to client %v: %v", client.RemoteAddr(), err)
				client.Close()
				delete(clients, client)
			}
			break
		}
	}
}
