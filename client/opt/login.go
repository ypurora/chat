package opt

import (
	"bytes"
	"chat/model"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Register(serverAddr string, user *model.User) bool {
	return sendCredentials(serverAddr, "/register", user)
}

func Login(serverAddr string, user *model.User) bool {
	return sendCredentials(serverAddr, "/login", user)
}

func sendCredentials(serverAddr, endpoint string, user *model.User) bool {
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Failed to marshal credentials: %v", err)
	}

	resp, err := http.Post("http://"+serverAddr+endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to send credentials: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Println("User not found. Please register first.")
		return false
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to %s: %s", strings.TrimPrefix(endpoint, "/"), resp.Status)
		return false
	}

	log.Printf("Successfully %sed", strings.TrimPrefix(endpoint, "/"))
	return true
}
