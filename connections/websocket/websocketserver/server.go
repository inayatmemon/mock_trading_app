package websocketserver

import (
	"9DotTechnology/trading_platform/server/streams"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"net/http"
)

func CreateWebSocketServer() error {
	// http.HandleFunc("/ws", handleWebSocket)
	streams.InitializeStreams()
	port := ":9000"
	logwrapper.Logger.Debugln("WebSocket Server started on port: ", port)
	go http.ListenAndServe(port, nil)
	return nil
}
