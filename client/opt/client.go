package opt

import (
	"bufio"
	"chat/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
)

func ReadMessages(conn *websocket.Conn, username string, done chan struct{}) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			close(done)
			return
		}

		var msg model.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			// If the message is not JSON, print it as plain text
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// Check if the message should be displayed based on the rooms
		shouldDisplay := false
		for _, room := range msg.Receiver {
			if room == "default" || room == currentRoom || msg.Type == "dm" {
				shouldDisplay = true
				break
			}
		}
		if shouldDisplay {
			log.Printf("%s: %s\n", msg.Sender, msg.Content)
		}
	}
}

func HandleUserInput(conn *websocket.Conn, username string, done chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages to send to the chat. Type 'exit' to quit.")
	fmt.Println("To send a DM, use the format: @username your message")
	fmt.Println("To create a room, use the command: /create room_name")
	fmt.Println("To join a room, use the command: /join room_name")
	fmt.Println("To leave a room, use the command: /leave room_name")

	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}

		var msg model.Message

		if strings.HasPrefix(text, "/create ") {
			roomName := strings.TrimSpace(strings.TrimPrefix(text, "/create "))
			msg = model.Message{
				Sender:  username,
				Content: roomName,
				Type:    "create_room",
			}
		} else if strings.HasPrefix(text, "/join ") {
			roomName := strings.TrimSpace(strings.TrimPrefix(text, "/join "))
			msg = model.Message{
				Sender:  username,
				Content: roomName,
				Type:    "join_room",
			}
			currentRoom = roomName
		} else if strings.HasPrefix(text, "/leave ") {
			roomName := strings.TrimSpace(strings.TrimPrefix(text, "/leave "))
			msg = model.Message{
				Sender:  username,
				Content: roomName,
				Type:    "leave_room",
			}
			currentRoom = "default"
		} else if strings.HasPrefix(text, "@") {
			parts := strings.SplitN(text, " ", 2)
			if len(parts) == 2 {
				receiver := strings.TrimPrefix(parts[0], "@")
				content := parts[1]
				msg = model.Message{
					Sender:   username,
					Content:  content,
					Receiver: []string{receiver},
					Type:     "dm",
				}
			} else {
				continue
			}
		} else {
			msg = createMessage(text, username)
			if currentRoom != "" {
				msg.Receiver = []string{currentRoom}
			} else {
				msg.Receiver = []string{"default"}
			}
		}

		err := conn.WriteJSON(&msg)
		if err != nil {
			log.Println("Write error:", err)
			close(done)
			break
		}
	}
}

// createMessage creates a Message struct from user input
func createMessage(text, username string) model.Message {
	var receiver []string
	content := text

	if strings.HasPrefix(text, "@") {
		parts := strings.SplitN(text, " ", 2)
		if len(parts) == 2 {
			receiver = append(receiver, strings.TrimPrefix(parts[0], "@"))
			content = parts[1]
		}
	}

	return model.Message{
		Sender:   username,
		Content:  content,
		Receiver: receiver,
	}
}

// ReadInput reads input from the standard input
func ReadInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
