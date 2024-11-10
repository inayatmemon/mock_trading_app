package tradehistory

import (
	"9DotTechnology/trading_platform/constants/db_constants"
	commonmodels "9DotTechnology/trading_platform/models/common"
	"9DotTechnology/trading_platform/models/v1/orders"
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/commons/db_helpers"
	"9DotTechnology/trading_platform/utils/localization"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTradeHistoryLogic(c *gin.Context, payload orders.TradeHistoryRequest) commonmodels.ApiResponse {
	lang := c.GetHeader("language")

	userDetails, err := commons.FetchUserDataFromAuthToken(c)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	getTradeHistoryQry := bson.M{
		"userId": userDetails.ID,
	}

	if payload.OrderId != "" {
		objId, _ := primitive.ObjectIDFromHex(payload.OrderId)
		getTradeHistoryQry["_id"] = objId
	}

	if payload.Symbol != "" {
		getTradeHistoryQry["symbol"] = payload.Symbol
	}

	history, err := db_helpers.GetManyDocuments[orders.Order](db_constants.OrdersCollection, getTradeHistoryQry)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	if len(history) == 0 {
		return commonmodels.ApiResponse{
			Code: http.StatusNoContent,
		}
	}

	return commonmodels.ApiResponse{
		Code:           http.StatusOK,
		Message:        localization.GetMessage(lang, "common.200", nil),
		Data:           history,
		TotalCount:     int64(len(history)),
		TotalPageCount: 1,
	}
}
