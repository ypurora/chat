package main

import (
	"chat/server/opt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialize user DB if it doesn't exist
	if _, err := os.Stat(opt.UserDB); os.IsNotExist(err) {
		if err := os.WriteFile(opt.UserDB, []byte("[]"), 0644); err != nil {
			log.Fatalf("Failed to initialize user DB: %v", err)
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(opt.OutputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Configure routes
	http.HandleFunc("/ws", opt.HandleConnections)
	http.HandleFunc("/register", opt.HandleRegister)
	http.HandleFunc("/login", opt.HandleLogin)

	// Start listening for incoming chat messages
	go opt.HandleMessages()

	// Start the server on localhost port 8080
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}
