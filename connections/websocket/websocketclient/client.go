package websocketclient

import (
	"9DotTechnology/trading_platform/constants/db_constants"
	marketDataModel "9DotTechnology/trading_platform/models/v1/market_data"
	marketdata "9DotTechnology/trading_platform/streams/market_data"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
)

var Conn *websocket.Conn

func createConnection(path string) (*websocket.Conn, error) {
	var err error
	if db_constants.WEB_SOCKET_HOST == "" || db_constants.WEB_SOCKET_SCHEME == "" {
		logwrapper.Logger.Debugln("CONFIGURATION MISSING FOR WEBSOCKET")
		return nil, errors.New("CONFIGURATION MISSING FOR WEBSOCKET")
	}

	serverURL := db_constants.WEB_SOCKET_SCHEME + "://" + db_constants.WEB_SOCKET_HOST + path
	logwrapper.Logger.Debugln("Connecting to web socket server at path: ", path)

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		logwrapper.Logger.Errorln("error in connecting to websocket server: ", err)
		return nil, err
	}

	logwrapper.Logger.Debugln("websocket server connected")

	return conn, nil
}

func SubscribeToBookTickerStream() error {
	conn, err := createConnection("/stream?streams=!bookTicker")
	if err != nil {
		logwrapper.Logger.Debugln("error in connecting to bookTicker: ", err)
		return err
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				logwrapper.Logger.Debugln("error in reading bookTicker:", err)
				conn.Close()
				break
			}

			var marketData marketDataModel.BookTickerMessage
			err = json.Unmarshal(message, &marketData)
			if err != nil {
				logwrapper.Logger.Debugln("error in unmarshalling market data: ", err)
				conn.Close()
				break
			}

			broadcastMessage := map[string]interface{}{
				"eventType":       marketData.Data.EventType,
				"updateID":        marketData.Data.UpdateID,
				"symbol":          marketData.Data.Symbol,
				"bidPrice":        marketData.Data.BidPrice,
				"bidQuantity":     marketData.Data.BidQuantity,
				"askPrice":        marketData.Data.AskPrice,
				"askQuantity":     marketData.Data.AskQuantity,
				"transactionTime": marketData.Data.TransactionTime,
				"eventTime":       marketData.Data.EventTime,
			}

			broadcastMessageBytes, _ := json.Marshal(broadcastMessage)

			// fmt.Printf("message: %v\n", message)

			// Forward the received Binance message to each connected client
			marketdata.BroadcastToClients(broadcastMessageBytes)
		}
	}()

	return nil
}

func SubscribeToLiveMarketPriceStream(symbol string) (*websocket.Conn, error) {
	conn, err := createConnection("/stream?streams=" + symbol + "@bookTicker")
	if err != nil {
		logwrapper.Logger.Debugln("error in connecting to bookTicker: ", err)
		return nil, err
	}

	return conn, nil
}
