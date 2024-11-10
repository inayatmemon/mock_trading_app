package position

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Position struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    string             `bson:"userId" json:"userId"`
	Symbol    string             `bson:"symbol" json:"symbol"`
	Volume    float64            `bson:"volume" json:"volume"`
	AvgPrice  float64            `bson:"avgPrice" json:"avgPrice"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type GetPositionPayload struct {
	Symbol string `form:"symbol" valid:"optional"`
}
