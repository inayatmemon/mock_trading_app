package orders

import (
	ordersLogic "9DotTechnology/trading_platform/buisnesslogics/v1/orders"
	"9DotTechnology/trading_platform/models/v1/orders"
	"9DotTechnology/trading_platform/utils/commons"

	"github.com/gin-gonic/gin"
)

// Api Handler for create order
func CreateOrder(c *gin.Context) {
	payload, err := commons.ParseJsonPayload[orders.OrderRequest](c)
	if err != nil {
		return
	}

	apiResponse := ordersLogic.CreateOrderLogic(c, payload)
	commons.SendApiResponse(c, apiResponse)
}

// Api Handler for delete order
func DeleteOrder(c *gin.Context) {
	payload, err := commons.ParseJsonPayload[orders.DeleteOrderRequest](c)
	if err != nil {
		return
	}

	apiResponse := ordersLogic.DeleteOrderLogic(c, payload)
	commons.SendApiResponse(c, apiResponse)
}
