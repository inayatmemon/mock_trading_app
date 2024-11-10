package streams

import (
	middleware "9DotTechnology/trading_platform/middlerware"
	marketdata "9DotTechnology/trading_platform/streams/market_data"
	"9DotTechnology/trading_platform/streams/orders"

	"net/http"
)

// function to live all the stream channels
func InitializeStreams() {
	http.HandleFunc("/market-data", middleware.ValidateStreamConnectionAuth(marketdata.GetMarketDataStream))
	http.HandleFunc("/order-create", middleware.ValidateStreamConnectionAuth(orders.CreateOrderStream))
	http.HandleFunc("/order-delete", middleware.ValidateStreamConnectionAuth(orders.DeleteOrderStream))
}
