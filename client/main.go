package main

import (
	"bufio"
	"chat/client/opt"
	"chat/model"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:8080"
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Chat App")

	var user model.User

	// Get user credentials
	fmt.Print("Enter username: ")
	user.Username = opt.ReadInput(reader)

	fmt.Print("Enter password: ")
	user.Password = opt.ReadInput(reader)

	// Choose to register or login
	fmt.Print("Do you want to (r)egister or (l)ogin? ")
	choice := strings.ToLower(strings.TrimSpace(opt.ReadInput(reader)))

	switch choice {
	case "r":
		if !opt.Register(serverAddr, &user) {
			return
		}
	case "l":
		if !opt.Login(serverAddr, &user) {
			fmt.Println("Login failed. Please make sure you have registered first.")
			return
		}
	default:
		fmt.Println("Invalid choice")
		return
	}

	conn, err := opt.ConnectToServer(serverAddr)
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	// Channel to signal end of application
	done := make(chan struct{})

	// Send initial auth message
	authMsg := model.Message{
		Sender:  user.Username,
		Content: user.Password, // Password sent for authentication, ideally use a token-based system for better security
		Type:    "auth",
	}

	if err := conn.WriteJSON(authMsg); err != nil {
		log.Fatal("Failed to send initial auth message:", err)
	}

	// join default room
	joinMsg := model.Message{
		Sender:  user.Username,
		Content: "default",
		Type:    "join_room",
	}
	if err := conn.WriteJSON(joinMsg); err != nil {
		log.Fatal("Failed to join default room:", err)
	}

	// Goroutine to read messages from the server
	go opt.ReadMessages(conn, user.Username, done)

	// Read user input from stdin and send to the server
	opt.HandleUserInput(conn, user.Username, done)

	// Wait for goroutine to finish
	<-done
	log.Println("Disconnected from server")
}
