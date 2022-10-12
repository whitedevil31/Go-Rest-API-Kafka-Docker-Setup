package config

import (
	"context"
	"fmt"
	"os"

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
		fmt.Println(err)
	}

	client = connectDB
	fmt.Println("Database connected")
	return client
}
func GetDB() *mongo.Client {
	if client == nil {
		return Connect()
	}
	return client
}
