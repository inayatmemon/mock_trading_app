package commons

import (
	"9DotTechnology/trading_platform/constants/db_constants"
	commonmodels "9DotTechnology/trading_platform/models/common"
	"9DotTechnology/trading_platform/models/v1/orders"
	"9DotTechnology/trading_platform/models/v1/position"
	"9DotTechnology/trading_platform/utils/auth"
	"9DotTechnology/trading_platform/utils/commons/db_helpers"
	"9DotTechnology/trading_platform/utils/localization"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Define an upgrader that will handle HTTP to WebSocket protocol upgrade
var WebSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Allow all origins (not recommended for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type invalidFieldsStruct struct {
	fieldName string
	customErr string
}

func ParseJsonPayload[T any](ctx *gin.Context) (T, error) {
	lang := ctx.GetHeader("language")
	response := GetDefaultApiResponse(400, lang)
	var requestData T
	err := ctx.ShouldBindJSON(&requestData)
	response.Message = checkForDataTypeErr(err, lang)
	if response.Message != "" {
		SendApiResponse(ctx, response)
		return requestData, err
	}

	requestData = TrimStringFields(requestData).(T)
	valid, err := govalidator.ValidateStruct(requestData)

	if err != nil || !valid {
		response.Message = validateFieldErrors(err, lang)
		SendApiResponse(ctx, response)
		return requestData, err
	}
	return requestData, nil
}

func ParseQueryPayload[T any](ctx *gin.Context) (T, error) {
	lang := ctx.GetHeader("language")
	response := GetDefaultApiResponse(400, lang)
	var requestData T
	err := ctx.ShouldBindQuery(&requestData)
	response.Message = checkForDataTypeErr(err, lang)
	if response.Message != "" {
		SendApiResponse(ctx, response)
		return requestData, err
	}

	requestData = TrimStringFields(requestData).(T)
	valid, err := govalidator.ValidateStruct(requestData)

	if err != nil || !valid {
		response.Message = validateFieldErrors(err, lang)
		SendApiResponse(ctx, response)
		return requestData, err
	}
	return requestData, nil
}

func checkForDataTypeErr(err error, lang string) (errMessage string) {
	var invalidType *json.UnmarshalTypeError
	if errors.As(err, &invalidType) {
		errMessage = localization.GetMessage(lang, "custom.invalid_data_type", map[string]interface{}{
			"Field":    invalidType.Field,
			"Expected": invalidType.Type,
			"Provided": invalidType.Value,
		})
		return
	}
	return
}

func validateFieldErrors(err error, lang string) (errMessage string) {
	requiredFields := []string{}
	invalidFields := []invalidFieldsStruct{}

	errors := govalidator.ErrorsByField(err)

	for field, err := range errors {
		if err == "non zero value required" {
			requiredFields = append(requiredFields, field)
		} else if strings.Contains(err, "does not validate as") {
			errWords := strings.Split(err, " ")
			invalidFields = append(invalidFields, invalidFieldsStruct{
				fieldName: field,
				customErr: errWords[len(errWords)-1],
			})
		}
	}

	if len(requiredFields) > 0 {
		field := "field"
		verb := "is"
		if len(requiredFields) > 1 {
			field = "fields"
			verb = "are"
		}
		errMessage = localization.GetMessage(lang, "common.400", map[string]interface{}{
			"Fields": strings.Join(requiredFields, ", "),
			"Field":  field,
			"Verb":   verb,
		})
		return
	} else if len(invalidFields) > 0 {
		// Set detailed error message based on field-specific validation error
		details := getInvalidFieldDetails(invalidFields[0].customErr)
		errMessage = localization.GetMessage(lang, "custom.invalid_fields", map[string]interface{}{
			"Field":   invalidFields[0].fieldName,
			"Details": details,
		})
		return
	}

	return
}

func getInvalidFieldDetails(customErr string) string {
	switch customErr {
	case "aadhar":
		return "it should have 12 digits and should not start with 0 & 1. Example: 234567890000"
	case "pan":
		return "it should have 10 letters; the first five are alphabets, then 4 digits, and an alphabet at the end. Example: ABCDE1234F"
	case "mobile":
		return "it should have 10 digits and start from 6 to 9. Example: 9123456780"
	case "objectid":
		return "it should follow the MongoDB ObjectId format. Example: 672fba4766589973dacdce75"
	case "ordertype":
		return "it should be 'buy' or 'sell'"
	case "ordercategory":
		return "it should be 'market', 'limit', or 'stop'"
	case "email":
		return "it should be in the format abcd@gmail.com"
	default:
		return ""
	}
}

func SendApiResponse(ctx *gin.Context, response commonmodels.ApiResponse) {
	if response.Code >= 200 && response.Code < 300 {
		ctx.JSON(response.Code, response)
	} else {
		ctx.AbortWithStatusJSON(response.Code, response)
	}
	return
}

func TrimStringFields(input interface{}) interface{} {
	// Create a new value from the input to avoid modifying the original
	val := reflect.ValueOf(input)

	// If the input is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure the input is a struct
	if val.Kind() != reflect.Struct {
		return input
	}

	// Create a new struct of the same type
	newStruct := reflect.New(val.Type()).Elem()

	// Iterate over each field in the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// fieldName := val.Type().Field(i).Name

		// Check if the field is a string
		if field.Kind() == reflect.String {
			// Trim whitespace and set the new value
			newValue := strings.TrimSpace(field.String())
			newStruct.Field(i).SetString(newValue)
		} else if field.Kind() == reflect.Struct {
			// Recursively handle nested structs
			newStruct.Field(i).Set(reflect.ValueOf(TrimStringFields(field.Interface())))
		} else if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			// Handle array of strings
			// Handle slices and arrays
			fmt.Println("yes it is array")
			newSlice := reflect.MakeSlice(field.Type(), field.Len(), field.Cap())
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.String {
					// Trim strings in the slice or array
					newSlice.Index(j).SetString(strings.TrimSpace(elem.String()))
				} else {
					// Copy other types as they are
					newSlice.Index(j).Set(elem)
				}
				// else if elem.Kind() == reflect.Struct {
				// 	// Recursively handle nested structs within the slice or array
				// 	newSlice.Index(j).Set(reflect.ValueOf(TrimStringFields(elem.Interface())))
				// } else if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
				// 	// Recursively handle nested slices or arrays
				// 	newSlice.Index(j).Set(reflect.ValueOf(TrimStringFields(elem.Interface())))
				// }
			}
			newStruct.Field(i).Set(newSlice)
		} else {
			// Copy other types as they are
			newStruct.Field(i).Set(field)
		}
	}

	return newStruct.Interface()
}

func GetDefaultApiResponse(statusCode int, lang string) commonmodels.ApiResponse {
	response := commonmodels.ApiResponse{
		Code:           http.StatusInternalServerError,
		Message:        localization.GetMessage(lang, "common.500", nil),
		Data:           nil,
		TotalCount:     0,
		TotalPageCount: 0,
	}

	if statusCode == 400 {
		response.Code = http.StatusBadRequest
		response.Message = localization.GetMessage(lang, "common.400", nil)
	} else if statusCode == 200 {
		response.Code = http.StatusOK
		response.Message = localization.GetMessage(lang, "common.200", nil)
	} else if statusCode == 201 {
		response.Code = http.StatusCreated
		response.Message = localization.GetMessage(lang, "common.201", nil)
	} else if statusCode == 401 {
		response.Code = http.StatusUnauthorized
		response.Message = localization.GetMessage(lang, "common.401", nil)
	} else if statusCode == 404 {
		response.Code = http.StatusNotFound
		response.Message = localization.GetMessage(lang, "common.404", nil)
	} else if statusCode == 403 {
		response.Code = http.StatusForbidden
		response.Message = localization.GetMessage(lang, "common.403", nil)
	}

	return response
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func FetchUserDataFromAuthToken(c *gin.Context) (auth.TokenMetadata, error) {
	var userData auth.TokenMetadata
	token := c.GetHeader("authorization")
	if token == "" {
		logwrapper.Logger.Debugln("token is empty cannot decode it.")
		return userData, errors.New("token is empty cannot decode it")
	}

	claims, err := auth.DecodeToken(token)
	if err != nil {
		logwrapper.Logger.Debugln("error in decoding token: ", err)
		return userData, err
	}

	md, err := json.Marshal(claims)
	if err != nil {
		logwrapper.Logger.Debugln("error in marshalling claims: ", err)
		return userData, err
	}

	err = json.Unmarshal(md, &userData)
	if err != nil {
		logwrapper.Logger.Debugln("error in unmarshalling claims: ", err)
		return userData, err
	}

	return userData, nil
}

func GetCurrentPosition(userID string, symbol string) (float64, error) {

	getPositionQry := bson.M{"userId": userID, "symbol": symbol}
	position, err := db_helpers.GetOneDocument[position.Position](db_constants.PositionsCollection, getPositionQry)

	if err == mongo.ErrNoDocuments {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return position.Volume, nil
}

func CalculateNewPosition(currentPosition float64, order orders.OrderRequest) float64 {
	if strings.ToLower(order.Type) == "buy" {
		return currentPosition + order.Volume
	}
	return currentPosition - order.Volume
}

// func ExecuteLimitAndStopOrders(marketData marketdata.BookTickerMessage) error {
// 	filter := bson.M{
// 		"symbol": marketData.Data.Symbol,
// 		"status": "pending",
// 	}

// 	count, err := db_helpers.GetCount(db_constants.OrdersCollection, filter)
// 	if err != nil {
// 		logwrapper.Logger.Debugln("error in getting pending order count: ", err)
// 		return err
// 	}

// 	if count > 0 {
// 		marketData, err := GetMarketData(symbol)
// 		if err != nil {
// 			logwrapper.Logger.Debugln("error in getting market data for pending orders: ", err)
// 			return err
// 		}

// 		orderPrice, err := GetOrderPrice(orderDe)

// 		updateData := bson.M{
// 			"status":    "executed",
// 			"updatedAt": time.Now(),
// 			"price":     "",
// 		}
// 	}
// }
