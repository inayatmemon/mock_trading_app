package db_constants

const (
	DB_URL         = "mongodb://localhost:27017"
	MONGO_DATABASE = "trading_platform"
	MongoTimeout   = 30

	WEB_SOCKET_HOST   = "fstream.binance.com"
	WEB_SOCKET_SCHEME = "wss"

	UsersCollection     = "users"
	PositionsCollection = "positions"
	OrdersCollection    = "orders"
)
