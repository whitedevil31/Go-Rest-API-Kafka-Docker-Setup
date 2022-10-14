package config

import (
	"context"

	"os"

	"github.com/whitedevil31/atlan-backend/producer-api/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx = context.TODO()

func Connect() *mongo.Client {
	db := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(db)
	connectDB, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.WarningLogger.Println("FAILED TO CONNECT TO DATABASE")
	}

	client = connectDB
	logger.InfoLogger.Println("CONNECTED  TO DATABASE")
	return client
}
func GetDB() *mongo.Client {
	if client == nil {
		return Connect()
	}
	return client
}
