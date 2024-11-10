package users

import (
	userslogic "9DotTechnology/trading_platform/buisnesslogics/v1/users"
	"9DotTechnology/trading_platform/models/v1/users"
	"9DotTechnology/trading_platform/utils/commons"

	"github.com/gin-gonic/gin"
)

func RegisterUserApi(c *gin.Context) {
	payload, err := commons.ParseJsonPayload[users.RegisterUserPayload](c)
	if err != nil {
		return
	}

	apiResponse := userslogic.RegisterUserLogic(c, payload)
	commons.SendApiResponse(c, apiResponse)
}

func LoginUserApi(c *gin.Context) {
	payload, err := commons.ParseJsonPayload[users.LoginUserPayload](c)
	if err != nil {
		return
	}

	apiResponse := userslogic.LoginUserLogic(c, payload)
	commons.SendApiResponse(c, apiResponse)
}
