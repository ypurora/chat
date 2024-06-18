package model

// Message struct
type Message struct {
	Sender   string   `json:"sender"`
	Content  string   `json:"content"`
	Receiver []string `json:"receiver,omitempty"` // Add optional Receiver field for multiple rooms
	Type     string   `json:"type,omitempty"`     // Message type for different message handling
}

// User struct for credentials
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // Hashed password
}
