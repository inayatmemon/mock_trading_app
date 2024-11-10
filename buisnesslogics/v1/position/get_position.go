package position

import (
	"9DotTechnology/trading_platform/constants/db_constants"
	commonmodels "9DotTechnology/trading_platform/models/common"
	"9DotTechnology/trading_platform/models/v1/position"
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/commons/db_helpers"
	"9DotTechnology/trading_platform/utils/localization"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPositionLogic(c *gin.Context, payload position.GetPositionPayload) commonmodels.ApiResponse {
	lang := c.GetHeader("language")

	// getting the user details from the token
	userDetails, err := commons.FetchUserDataFromAuthToken(c)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	// applying the position filter
	getPositionQry := bson.M{
		"userId": userDetails.ID,
	}

	if payload.Symbol != "" {
		getPositionQry["symbol"] = payload.Symbol
	}

	// getting the position from the db
	position, err := db_helpers.GetManyDocuments[position.Position](db_constants.PositionsCollection, getPositionQry)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	// sending the response
	if len(position) == 0 {
		return commonmodels.ApiResponse{
			Code: http.StatusNoContent,
		}
	}

	return commonmodels.ApiResponse{
		Code:           http.StatusOK,
		Message:        localization.GetMessage(lang, "common.200", nil),
		Data:           position,
		TotalCount:     int64(len(position)),
		TotalPageCount: 1,
	}
}
