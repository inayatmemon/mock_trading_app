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
func DeleteOrderStream(w http.ResponseWriter, r *http.Request) {
	conn, err := commons.WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logwrapper.Logger.Debugln("Upgrade error:", err)
		return
	}

	defer conn.Close()

	// rading the delete order payload
	go func() {
		for {
			_, deleteOrderData, err := conn.ReadMessage()
			if err != nil {
				logwrapper.Logger.Debugln("error in getting payload for delete order stream: ", err)
				return
			}

			ctx := gin.Context{}
			ctx.Request, err = http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(deleteOrderData))
			if err != nil {
				logwrapper.Logger.Debugln("error in creating new request from delete order stream: ", err)
				return
			}
			ctx.Request.Header.Set("authorization", r.Header.Get("authorization"))
			ctx.Request.Header.Set("laguage", "en")
			ctx.Request.Header.Set("content-type", r.Header.Get("application/json"))

			// validate payload
			payload, err := commons.ParseJsonPayload[orders.DeleteOrderRequest](&ctx)
			if err != nil {
				logwrapper.Logger.Debugln("invalid payload for delete order stream: ", err)
				return
			}

			deleteOrderResp := orderLogics.DeleteOrderLogic(&ctx, payload)
			deleteOrderRespBytes, _ := json.Marshal(deleteOrderResp)
			conn.WriteMessage(1, deleteOrderRespBytes)
		}
	}()

	// Keep connection open for client to receive broadcast messages
	select {}
}
