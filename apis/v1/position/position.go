package position

import (
	positionLogic "9DotTechnology/trading_platform/buisnesslogics/v1/position"
	"9DotTechnology/trading_platform/models/v1/position"
	"9DotTechnology/trading_platform/utils/commons"

	"github.com/gin-gonic/gin"
)

func GetPosition(c *gin.Context) {
	payload, err := commons.ParseQueryPayload[position.GetPositionPayload](c)
	if err != nil {
		return
	}

	apiResponse := positionLogic.GetPositionLogic(c, payload)
	commons.SendApiResponse(c, apiResponse)
}
