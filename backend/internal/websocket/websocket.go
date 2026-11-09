// Package websocket keeps a set of authenticated clients and broadcasts
// small task-change events to them so task lists can refresh.
package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"task-management-system/internal/auth"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

type client struct {
	conn *websocket.Conn
	mu   sync.Mutex // gorilla connections allow one concurrent writer
}

func (c *client) write(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, data)
}

var (
	clientsMu sync.Mutex
	clients   = map[*client]struct{}{}
)

// Event is the JSON message sent to clients.
type Event struct {
	Type   string `json:"type"`
	TaskID uint   `json:"taskId"`
}

// Handler returns the /ws handler. Browsers must connect from one of the
// allowed origins and pass a valid JWT as ?token=... (browsers cannot set an
// Authorization header on WebSocket requests).
func Handler(allowedOrigins []string) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := strings.TrimRight(r.Header.Get("Origin"), "/")
			if origin == "" {
				return true // non-browser client; still needs a valid token
			}
			for _, allowed := range allowedOrigins {
				if allowed == "*" || strings.EqualFold(origin, allowed) {
					return true
				}
			}
			return false
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if _, _, err := auth.VerifyToken(r.URL.Query().Get("token")); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		c := &client{conn: conn}

		clientsMu.Lock()
		clients[c] = struct{}{}
		clientsMu.Unlock()

		defer func() {
			clientsMu.Lock()
			delete(clients, c)
			clientsMu.Unlock()
			conn.Close()
		}()

		conn.SetReadLimit(1024)
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(pongWait))
		})

		done := make(chan struct{})
		defer close(done)
		go func() {
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := c.write(websocket.PingMessage, nil); err != nil {
						return
					}
				case <-done:
					return
				}
			}
		}()

		// Clients only listen; incoming messages are read and discarded so
		// that control frames (pong/close) are processed.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}
}

// BroadcastTaskEvent notifies every connected client that a task changed.
// Only the task ID is sent; clients refetch through the authorised API.
func BroadcastTaskEvent(eventType string, taskID uint) {
	data, err := json.Marshal(Event{Type: eventType, TaskID: taskID})
	if err != nil {
		return
	}

	clientsMu.Lock()
	targets := make([]*client, 0, len(clients))
	for c := range clients {
		targets = append(targets, c)
	}
	clientsMu.Unlock()

	for _, c := range targets {
		if err := c.write(websocket.TextMessage, data); err != nil {
			c.conn.Close()
		}
	}
}
