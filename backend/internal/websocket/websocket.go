package websocket

import (
    "github.com/gorilla/websocket"
    "net/http"
    "sync"
    "log"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

// ✅ WebSocket Handler
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade failed:", err)
        return
    }
    defer conn.Close()

    log.Println("WebSocket connected:", conn.RemoteAddr())

    mutex.Lock()
    clients[conn] = true
    mutex.Unlock()

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println("WebSocket read error:", err)
            mutex.Lock()
            delete(clients, conn)
            mutex.Unlock()
            return
        }

        log.Println("WebSocket received:", string(p))
        BroadcastMessage(string(p)) // ✅ Broadcast to all clients

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println("WebSocket write error:", err)
            return
        }
    }
}

// ✅ Function to broadcast messages to all clients
func BroadcastMessage(message string) {
    mutex.Lock()
    defer mutex.Unlock()

    for client := range clients {
        err := client.WriteMessage(websocket.TextMessage, []byte(message))
        if err != nil {
            client.Close()
            delete(clients, client)
        }
    }
}
