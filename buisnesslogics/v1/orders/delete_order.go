package orders

import (
	"9DotTechnology/trading_platform/constants/db_constants"
	commonmodels "9DotTechnology/trading_platform/models/common"
	"9DotTechnology/trading_platform/models/v1/orders"
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/commons/db_helpers"
	"9DotTechnology/trading_platform/utils/localization"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// buisness logic function for deleting the order if it has pending status
func DeleteOrderLogic(c *gin.Context, payload orders.DeleteOrderRequest) commonmodels.ApiResponse {
	lang := c.GetHeader("language")

	objId, _ := primitive.ObjectIDFromHex(payload.OrderId)

	// applying the order id filter
	filter := bson.M{
		"_id": objId,
	}

	// getting the order details from the db
	orderDetails, err := db_helpers.GetOneDocument[orders.Order](db_constants.OrdersCollection, filter)
	if err == mongo.ErrNoDocuments {
		return commons.GetDefaultApiResponse(http.StatusNotFound, lang)
	}

	// sending the message to user if the status of the order is executed
	if orderDetails.Status == "executed" {
		return commonmodels.ApiResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: localization.GetMessage(lang, "orders.already_executed", nil),
		}
	}

	// sending the message to user if the status of the order is deleted
	if orderDetails.Status == "deleted" {
		return commonmodels.ApiResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: localization.GetMessage(lang, "orders.already_deleted", nil),
		}
	}

	// updating the status as deleted in db
	updateData := bson.M{
		"status":    "deleted",
		"updatedAt": time.Now(),
	}

	err = db_helpers.UpdateOneDocumentWithID(db_constants.OrdersCollection, orderDetails.ID, updateData)
	if err != nil {
		logwrapper.Logger.Debugln("error in deleting the order: ", err)
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	return commons.GetDefaultApiResponse(http.StatusOK, lang)
}
