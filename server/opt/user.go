package opt

import (
	"chat/model"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func handleAuth(ws *websocket.Conn, msg model.Message) {
	if isValidUser(msg.Sender, msg.Content) {
		clientsMux.Lock()
		clients[ws] = msg.Sender
		clientsMux.Unlock()
		log.Printf("Client authenticated: %s", msg.Sender)
	} else {
		log.Printf("Authentication failed for: %s", msg.Sender)
		clientsMux.Lock()
		delete(clients, ws)
		clientsMux.Unlock()
		ws.WriteMessage(websocket.TextMessage, []byte("Authentication failed"))
	}
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user.Password = hashPassword(user.Password)

	users, err := loadUsers()
	if err != nil {
		log.Printf("Error loading users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		if u.Username == user.Username {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
	}

	users = append(users, user)
	if err := saveUsers(users); err != nil {
		log.Printf("Error saving users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	users, err := loadUsers()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		if u.Username == user.Username {
			if comparePasswords(u.Password, user.Password) {
				w.WriteHeader(http.StatusOK)
				return
			} else {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func loadUsers() ([]model.User, error) {
	data, err := ioutil.ReadFile(UserDB)
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func saveUsers(users []model.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(UserDB, data, 0644)
}

func hashPassword(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	return hashedPwd == hashPassword(plainPwd)
}

func isValidUser(username, password string) bool {
	users, err := loadUsers()
	if err != nil {
		log.Printf("Error loading users: %v", err)
		return false
	}

	for _, u := range users {
		if u.Username == username {
			return comparePasswords(u.Password, password)
		}
	}

	return false
}
