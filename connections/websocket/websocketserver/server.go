package websocketserver

import (
	"9DotTechnology/trading_platform/server/streams"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"net/http"
)

// function to create local websocket server
func CreateWebSocketServer() error {
	// initializing the local streams endpoints
	streams.InitializeStreams()
	port := ":9000"
	logwrapper.Logger.Debugln("WebSocket Server started on port: ", port)
	go http.ListenAndServe(port, nil)
	return nil
}
