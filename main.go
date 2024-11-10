package main

import (
	"9DotTechnology/trading_platform/connections/mongodb"
	"9DotTechnology/trading_platform/connections/websocket/websocketclient"
	"9DotTechnology/trading_platform/connections/websocket/websocketserver"

	"9DotTechnology/trading_platform/utils/localization"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"9DotTechnology/trading_platform/utils/validator"

	"9DotTechnology/trading_platform/server"
	"strings"
)

func main() {

	validator.Init()

	Logger := logwrapper.NewLogger()
	localization.LoadBundle()

	Logger.Infoln(strings.Repeat("~", 50))
	mongoErr := mongodb.NewConnection()
	if mongoErr != nil {
		Logger.Fatal("Error connecting to MongoDB : ", mongoErr)
	}

	Logger.Infoln(strings.Repeat("~", 50))
	wsErr := websocketserver.CreateWebSocketServer()
	if wsErr != nil {
		Logger.Fatal("Error connecting to WebSocket : ", wsErr)
	}

	cliErr := websocketclient.SubscribeToBookTickerStream()
	if cliErr != nil {
		Logger.Fatal("Error connecting to WebSocket Binance Server: ", cliErr)
	}

	server.StartServer()

}
