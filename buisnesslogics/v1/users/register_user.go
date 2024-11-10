package users

import (
	"9DotTechnology/trading_platform/constants/db_constants"
	commonmodels "9DotTechnology/trading_platform/models/common"
	"9DotTechnology/trading_platform/models/v1/users"
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/commons/db_helpers"
	"9DotTechnology/trading_platform/utils/localization"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterUserLogic(c *gin.Context, payload users.RegisterUserPayload) commonmodels.ApiResponse {
	lang := c.GetHeader("language")
	conflictQry := bson.M{
		"$or": []bson.M{
			{
				"countryCode": payload.CountryCode,
				"mobile":      payload.Mobile,
			},
			{
				"email": payload.Email,
			},
			{
				"aadharNumber": payload.AadharNumber,
			},
			{
				"panNumber": payload.PanNumber,
			},
		},
	}

	count, err := db_helpers.GetCount(db_constants.UsersCollection, conflictQry)
	if err != nil {
		logwrapper.Logger.Debugln("error in users count qry: ", err)
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	if count > 0 {
		return commonmodels.ApiResponse{
			Code:    http.StatusConflict,
			Message: localization.GetMessage(lang, "users.conflict", nil),
		}
	}

	password, err := commons.HashPassword(payload.Password)
	if err != nil {
		logwrapper.Logger.Debugln("error in password hash: ", err)
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	userData := users.Users{
		ID:            primitive.NewObjectID(),
		FirstName:     payload.FirstName,
		LastName:      payload.LastName,
		Email:         payload.Email,
		Password:      password,
		AadharNumber:  payload.AadharNumber,
		PanNumber:     payload.PanNumber,
		CountryCode:   payload.CountryCode,
		Mobile:        payload.Mobile,
		Status:        "ACTIVE",
		CreatedOn:     time.Now().Unix(),
		LastUpdatedOn: time.Now().Unix(),
	}

	err = db_helpers.InsertOneDocument(db_constants.UsersCollection, userData)
	if err != nil {
		logwrapper.Logger.Debugln("error in registering the user: ", err)
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	return commons.GetDefaultApiResponse(http.StatusCreated, lang)
}
