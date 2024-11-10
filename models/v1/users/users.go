package users

import "go.mongodb.org/mongo-driver/bson/primitive"

// users collection db model
type Users struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	FirstName     string             `json:"firstName" bson:"firstName"`
	LastName      string             `json:"lastName" bson:"lastName"`
	CountryCode   string             `json:"countryCode" bson:"countryCode"`
	Mobile        string             `json:"mobile" bson:"mobile"`
	Email         string             `json:"email" bson:"email"`
	Password      string             `json:"password" bson:"password"`
	AadharNumber  string             `json:"aadharNumber" bson:"aadharNumber"`
	PanNumber     string             `json:"panNumber" bson:"panNumber"`
	Status        string             `json:"status" bson:"status"`
	CreatedOn     int64              `json:"createdOn" bson:"createdOn"`
	LastUpdatedOn int64              `json:"lastUpdatedOn" bson:"lastUpdatedOn"`
}

// register user request model
type RegisterUserPayload struct {
	FirstName    string `json:"firstName" valid:"required"`
	LastName     string `json:"lastName" valid:"optional"`
	Password     string `json:"password" valid:"required"`
	CountryCode  string `json:"countryCode" valid:"required"`
	Mobile       string `json:"mobile" valid:"required,mobile"`
	Email        string `json:"email" valid:"required,email"`
	AadharNumber string `json:"aadharNumber" valid:"required,aadhar"`
	PanNumber    string `json:"panNumber" valid:"required,pan"`
}

// login user request model
type LoginUserPayload struct {
	Password string `json:"password" valid:"required"`
	Email    string `json:"email" valid:"required,email"`
}
