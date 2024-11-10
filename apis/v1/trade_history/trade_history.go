package tradehistory

import (
	tradeHistoryLogic "9DotTechnology/trading_platform/buisnesslogics/v1/trade_history"
	"9DotTechnology/trading_platform/models/v1/orders"
	"9DotTechnology/trading_platform/utils/commons"

	"github.com/gin-gonic/gin"
)

// Api Handler for getting trade history
func GetTradeHistoryApi(c *gin.Context) {
	payload, err := commons.ParseQueryPayload[orders.TradeHistoryRequest](c)
	if err != nil {
		return
	}

	apiResponse := tradeHistoryLogic.GetTradeHistoryLogic(c, payload)
	commons.SendApiResponse(c, apiResponse)
}
