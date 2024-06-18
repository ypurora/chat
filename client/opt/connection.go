package opt

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
)

func ConnectToServer(serverAddr string) (*websocket.Conn, error) {
	var conn *websocket.Conn
	var err error

	maxRetries := 5
	baseDelay := time.Second
	maxDelay := 30 * time.Second
	for i := 0; i < maxRetries; i++ {
		conn, _, err = websocket.DefaultDialer.Dial("ws://"+serverAddr+"/ws", nil)
		if err == nil {
			break
		}

		delay := calculateBackoff(i, baseDelay, maxDelay)
		log.Printf("Failed to connect, retrying in %v... (%d/%d)", delay, i+1, maxRetries)
		time.Sleep(delay)
	}

	return conn, err
}

// calculateBackoff returns the delay duration for the given retry count
func calculateBackoff(retryCount int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := baseDelay * time.Duration(1<<uint(retryCount))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay + time.Duration(rand.Intn(int(delay/2))) // Add random jitter
}
