package orders

import (
	orderLogics "9DotTechnology/trading_platform/buisnesslogics/v1/orders"
	"9DotTechnology/trading_platform/models/v1/orders"
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Initialize WebSocket and start broadcaster
func CreateOrderStream(w http.ResponseWriter, r *http.Request) {
	conn, err := commons.WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logwrapper.Logger.Debugln("Upgrade error:", err)
		return
	}

	defer conn.Close()
	go func() {
		for {
			_, createOrderData, err := conn.ReadMessage()
			if err != nil {
				logwrapper.Logger.Debugln("error in getting payload for create order stream: ", err)
				return
			}

			ctx := gin.Context{}
			ctx.Request, err = http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(createOrderData))
			if err != nil {
				logwrapper.Logger.Debugln("error in creating new request from create order stream: ", err)
				return
			}
			ctx.Request.Header.Set("authorization", r.Header.Get("authorization"))
			ctx.Request.Header.Set("laguage", "en")
			ctx.Request.Header.Set("content-type", r.Header.Get("application/json"))

			payload, err := commons.ParseJsonPayload[orders.OrderRequest](&ctx)
			if err != nil {
				logwrapper.Logger.Debugln("invalid payload for create order stream: ", err)
				return
			}

			createOrderResp := orderLogics.CreateOrderLogic(&ctx, payload)
			createOrderRespBytes, _ := json.Marshal(createOrderResp)
			conn.WriteMessage(1, createOrderRespBytes)
		}
	}()

	// Keep connection open for client to receive broadcast messages
	select {}
}
