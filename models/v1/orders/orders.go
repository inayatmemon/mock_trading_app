package orders

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// orders collection db model
type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"userId" json:"userId"`
	Symbol          string             `bson:"symbol" json:"symbol"`
	Volume          float64            `bson:"volume" json:"volume"`
	Type            string             `bson:"type" json:"type"`
	OrderCategory   string             `bson:"orderCategory" json:"orderCategory"`
	PricePerVoulume float64            `bson:"pricePerVoulume" json:"pricePerVoulume"`
	LimitPrice      float64            `bson:"limitPrice" json:"limitPrice"`
	StopPrice       float64            `bson:"stopPrice" json:"stopPrice"`
	TotalPrice      float64            `bson:"totalPrice" json:"totalPrice"`
	Status          string             `bson:"status" json:"status"` // "executed" or "pending"
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// create order payload model
type OrderRequest struct {
	Symbol        string  `json:"symbol" valid:"required"`
	Volume        float64 `json:"volume" valid:"required"`
	Type          string  `json:"type" valid:"required,ordertype"`              // "buy" or "sell"
	OrderCategory string  `json:"orderCategory" valid:"required,ordercategory"` // "market", "limit", "stop"
	LimitPrice    float64 `json:"limitPrice" valid:"optional"`                  // Required for limit orders
	StopPrice     float64 `json:"stopPrice" valid:"optional"`                   // Required for stop orders
}

// get trade history api request model
type TradeHistoryRequest struct {
	OrderId string `form:"orderId" valid:"optional"`
	Symbol  string `form:"symbol" valid:"optional"`
}

// delete order from db request model
type DeleteOrderRequest struct {
	OrderId string `form:"orderId" valid:"required,objectid"`
}
