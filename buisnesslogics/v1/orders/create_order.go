package orders

import (
	"9DotTechnology/trading_platform/connections/websocket/websocketclient"
	"9DotTechnology/trading_platform/constants/db_constants"
	commonmodels "9DotTechnology/trading_platform/models/common"
	marketdata "9DotTechnology/trading_platform/models/v1/market_data"
	"9DotTechnology/trading_platform/models/v1/orders"
	positionModel "9DotTechnology/trading_platform/models/v1/position"
	"9DotTechnology/trading_platform/utils/commons"
	"9DotTechnology/trading_platform/utils/commons/db_helpers"
	"9DotTechnology/trading_platform/utils/localization"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrderLogic(c *gin.Context, payload orders.OrderRequest) commonmodels.ApiResponse {
	lang := c.GetHeader("language")

	payload.Type = strings.ToUpper(payload.Type)
	payload.OrderCategory = strings.ToUpper(payload.OrderCategory)
	if payload.OrderCategory != "MARKET" {
		if payload.OrderCategory == "STOP" && payload.StopPrice == 0 {
			return commonmodels.ApiResponse{
				Code: http.StatusBadRequest,
				Message: localization.GetMessage(lang, "orders.price", map[string]interface{}{
					"Type": "stopPrice",
				}),
			}
		}
		if payload.OrderCategory == "LIMIT" && payload.LimitPrice == 0 {
			return commonmodels.ApiResponse{
				Code: http.StatusBadRequest,
				Message: localization.GetMessage(lang, "orders.price", map[string]interface{}{
					"Type": "limitPrice",
				}),
			}
		}
	}

	userDetails, err := commons.FetchUserDataFromAuthToken(c)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	// Get market data from Binance WebSocket
	marketData, err := GetMarketData(payload.Symbol)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	// Check position limit (10 lots)
	currentPosition, err := commons.GetCurrentPosition(userDetails.ID, payload.Symbol)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	newPosition := commons.CalculateNewPosition(currentPosition, payload)
	if newPosition > 10 && strings.ToLower(payload.Type) == "buy" {
		return commonmodels.ApiResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: localization.GetMessage(lang, "orders.exceeds_10_lots", nil),
		}
	} else if newPosition < 0 && strings.ToLower(payload.Type) == "sell" {
		return commonmodels.ApiResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: localization.GetMessage(lang, "orders.negative", nil),
		}
	}

	orderPrice, err := GetOrderPrice(payload.Type, &marketData)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	orderId := primitive.NewObjectID()

	// Create order
	orderData := orders.Order{
		ID:            orderId,
		UserID:        userDetails.ID,
		Symbol:        payload.Symbol,
		Volume:        payload.Volume,
		Type:          payload.Type,
		OrderCategory: payload.OrderCategory,
		Status:        determineOrderStatus(payload.OrderCategory),
		CreatedAt:     time.Now(),
	}

	// Set prices based on order category
	switch orderData.OrderCategory {
	case "MARKET":
		orderData.PricePerVoulume = orderPrice
		orderData.TotalPrice = orderPrice * payload.Volume
	case "LIMIT":
		orderData.LimitPrice = payload.LimitPrice
	case "STOP":
		orderData.StopPrice = payload.StopPrice
	}

	err = db_helpers.InsertOneDocument(db_constants.OrdersCollection, orderData)
	if err != nil {
		return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
	}

	if payload.OrderCategory == "MARKET" {
		err = updatePosition(userDetails.ID, orderData)
		if err != nil {
			return commons.GetDefaultApiResponse(http.StatusInternalServerError, lang)
		}
	}

	return commonmodels.ApiResponse{
		Code:    http.StatusOK,
		Message: localization.GetMessage(lang, "common.200", nil),
		Data: map[string]interface{}{
			"orderId":   orderId.Hex(),
			"price":     orderPrice,
			"status":    determineOrderStatus(orderData.OrderCategory),
			"qauntity":  payload.Volume,
			"timestamp": time.Now().UnixMilli(),
		},
	}
}

func determineOrderStatus(orderType string) string {
	if orderType == "MARKET" {
		return "executed"
	}
	return "pending"
}

func updatePosition(userID string, order orders.Order) error {

	filter := bson.M{
		"userId": userID,
		"symbol": order.Symbol,
	}

	position, err := db_helpers.GetOneDocument[positionModel.Position](db_constants.PositionsCollection, filter)
	if err == mongo.ErrNoDocuments {
		// Create new position
		position = positionModel.Position{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Symbol:    order.Symbol,
			Volume:    getPositionVolume(order),
			AvgPrice:  order.PricePerVoulume,
			UpdatedAt: time.Now(),
		}
		err = db_helpers.InsertOneDocument(db_constants.PositionsCollection, position)
		return err
	}
	if err != nil {
		logwrapper.Logger.Debugln("error in getting positions from db: ", err)
		return err
	}

	// Update existing position
	newVolume := position.Volume + getPositionVolume(order)
	var newAvgPrice float64

	if newVolume != 0 {
		if position.Volume == 0 {
			newAvgPrice = order.PricePerVoulume
		} else {
			totalValue := position.AvgPrice*position.Volume + order.PricePerVoulume*getPositionVolume(order)
			newAvgPrice = totalValue / newVolume
		}
	} else {
		newAvgPrice = 0
	}

	update := bson.M{
		"$set": bson.M{
			"volume":    newVolume,
			"avgPrice":  newAvgPrice,
			"updatedAt": time.Now(),
		},
	}

	err = db_helpers.UpdateOneDocumentWithData(db_constants.PositionsCollection, filter, update)
	return err
}

func getPositionVolume(order orders.Order) float64 {
	if strings.ToLower(order.Type) == "buy" {
		return order.Volume
	}
	return -order.Volume
}

func GetMarketData(symbol string) (marketdata.BookTickerData, error) {
	var marketData marketdata.BookTickerData
	var marketDataMessage marketdata.BookTickerMessage
	conn, err := websocketclient.SubscribeToLiveMarketPriceStream(symbol)
	if err != nil {
		logwrapper.Logger.Debugln("error in connecting to market price stream: ", err)
		return marketData, err
	}

	defer conn.Close()

	_, message, err := conn.ReadMessage()
	if err != nil {
		logwrapper.Logger.Debugln("error in reading message from market price stream: ", err)
		return marketData, err
	}

	// fmt.Printf("string(message): %v\n", string(message))

	err = json.Unmarshal(message, &marketDataMessage)
	if err != nil {
		logwrapper.Logger.Debugln("error in unmarshalling message from market price stream: ", err)
		return marketData, err
	}

	marketData = marketDataMessage.Data
	return marketData, nil
}

func GetOrderPrice(orderType string, marketData *marketdata.BookTickerData) (float64, error) {
	if strings.ToLower(orderType) == "buy" {
		fmt.Printf("marketData.AskPrice: %v\n", marketData.AskPrice)
		askPrice, err := strconv.ParseFloat(marketData.AskPrice, 64)
		if err != nil {
			logwrapper.Logger.Debugln("error in converting ask price to float: ", err)
			return 0, err
		}
		return askPrice, nil
	}

	bidPrice, err := strconv.ParseFloat(marketData.BidPrice, 64)
	if err != nil {
		logwrapper.Logger.Debugln("error in converting bid price to float: ", err)
		return 0, err
	}
	return bidPrice, nil
}
