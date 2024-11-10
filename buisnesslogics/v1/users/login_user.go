package users

import (
	"9DotTechnology/trading_platform/constants/common_constants"
	"9DotTechnology/trading_platform/constants/db_constants"
	commonmodels "9DotTechnology/trading_platform/models/common"
	"9DotTechnology/trading_platform/models/v1/users"
	"9DotTechnology/trading_platform/utils/auth"
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/commons/db_helpers"
	"9DotTechnology/trading_platform/utils/localization"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func LoginUserLogic(c *gin.Context, payload users.LoginUserPayload) commonmodels.ApiResponse {
	lang := c.GetHeader("language")
	existingUserQry := bson.M{
		"email": payload.Email,
	}

	count, err := db_helpers.GetCount(db_constants.UsersCollection, existingUserQry)
	if err != nil {
		logwrapper.Logger.Debugln("error in users count qry: ", err)
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	if count == 0 {
		return commonmodels.ApiResponse{
			Code:    http.StatusNotFound,
			Message: localization.GetMessage(lang, "users.not_registered", map[string]interface{}{"Email": payload.Email}),
		}
	}

	userData, err := db_helpers.GetOneDocument[users.Users](db_constants.UsersCollection, existingUserQry)
	if err != nil {
		logwrapper.Logger.Debugln("error in getting user from db: ", err)
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	if userData.Status != "ACTIVE" {
		return commonmodels.ApiResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: localization.GetMessage(lang, "users.inactive_account", nil),
		}
	}

	passwordMatched := commons.CheckPasswordHash(payload.Password, userData.Password)
	if !passwordMatched {
		return commonmodels.ApiResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: localization.GetMessage(lang, "users.password", nil),
		}
	}

	userMetadata := auth.TokenMetadata{
		Email:       userData.Email,
		ID:          userData.ID.Hex(),
		Mobile:      userData.Mobile,
		CountryCode: userData.CountryCode,
	}

	token, err := auth.GenerateToken(userMetadata)
	if err != nil {
		logwrapper.Logger.Debugln("error in generating the jwt token: ", err)
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	return commonmodels.ApiResponse{
		Code:    http.StatusOK,
		Message: localization.GetMessage(lang, "common.200", nil),
		Data: map[string]interface{}{
			"userId": userData.ID.Hex(),
			"email":  userData.Email,
			"tokenDetails": map[string]interface{}{
				"token":     token,
				"expiresAt": time.Now().Unix() + (common_constants.TOKEN_EXPIRY_HOURS * 3600),
			},
		},
	}
}
