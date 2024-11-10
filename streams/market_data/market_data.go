package marketdata

import (
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients = make(map[*websocket.Conn]bool) // connected clients
	// broadcast = make(chan []byte)              // broadcast channel
	mutex sync.Mutex // to synchronize access to the clients map
)

// Initialize WebSocket and start broadcaster
func GetMarketDataStream(w http.ResponseWriter, r *http.Request) {
	conn, err := commons.WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logwrapper.Logger.Debugln("Upgrade error:", err)
		return
	}

	// Register new client
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Remove client on disconnect
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
		logwrapper.Logger.Debugln("Client disconnected")
	}()

	// Keep connection open for client to receive broadcast messages
	select {}
}

// BroadcastToClients forwards messages to all active client connections
func BroadcastToClients(message []byte) {
	// Lock the clients map to ensure safe access
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logwrapper.Logger.Debugln("Error broadcasting message to client:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
