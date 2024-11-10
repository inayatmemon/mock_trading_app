package v1

import (
	"9DotTechnology/trading_platform/apis/v1/orders"
	"9DotTechnology/trading_platform/apis/v1/position"
	tradehistory "9DotTechnology/trading_platform/apis/v1/trade_history"
	"9DotTechnology/trading_platform/apis/v1/users"
	middleware "9DotTechnology/trading_platform/middlerware"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(router *gin.RouterGroup) {
	router.POST("/register", middleware.ValidateHeaders(), users.RegisterUserApi)
	router.POST("/login", middleware.ValidateHeaders(), users.LoginUserApi)
	router.POST("/order", middleware.ValidateAuthHeaders(), middleware.ValidateAuthToken(), orders.CreateOrder)
	router.DELETE("/order", middleware.ValidateAuthHeaders(), middleware.ValidateAuthToken(), orders.DeleteOrder)

	router.GET("/trade-history", middleware.ValidateAuthHeaders(), middleware.ValidateAuthToken(), tradehistory.GetTradeHistoryApi)
	router.GET("/position", middleware.ValidateAuthHeaders(), middleware.ValidateAuthToken(), position.GetPosition)
}
