package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

func init() {

}

func Connect() mongo.Client {
	MONGO_URI := os.Getenv("DB_URI")

	clientOption := options.Client().ApplyURI(MONGO_URI)

	client, err := mongo.Connect(context.Background(), clientOption)

	if err != nil {
		log.Fatal("Error connecting to database")

	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("Error connecting to database")

	}
	Collection = client.Database("golang_db").Collection("todos")
	fmt.Println("Connected to MongoDB database")
	return *client

}
