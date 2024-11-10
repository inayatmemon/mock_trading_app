package mongodb

import (
	"9DotTechnology/trading_platform/constants/db_constants"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection - Connection structure
type Connection struct {
	Conn     *mongo.Client
	ConnDB   *mongo.Database
	Database string
}

// Client - MongoDB Connection Golbal variable
var Client *Connection

// NewConnection - new connection of mongodb
func NewConnection() error {
	if db_constants.DB_URL == "" || db_constants.MONGO_DATABASE == "" {
		return fmt.Errorf("CONFIGURATION IS MISSING FOR MONGODB")
	}

	mongoClient := &Connection{
		Conn:     nil,
		ConnDB:   nil,
		Database: db_constants.MONGO_DATABASE,
	}
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(db_constants.DB_URL)

	var err error
	mongoClient.Conn, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("MONGO ERROR: %s", err)
	}
	err = mongoClient.Conn.Ping(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("MONGO PING ERROR: %s", err)
	}

	logwrapper.Logger.Infoln("Connected to MONGODB")

	mongoClient.ConnDB = mongoClient.Conn.Database(mongoClient.Database)
	Client = mongoClient

	return nil
}

// GetCollection - Helper Functions
func GetCollection(collectionName string) *mongo.Collection {
	return Client.ConnDB.Collection(collectionName)
}

// DbContext - Helper Functions
func DbContext(i time.Duration) (context.Context, context.CancelFunc) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), i*time.Second)
	return ctx, ctxCancel
}
