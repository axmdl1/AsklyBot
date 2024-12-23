package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Монгоға қосылған варейбл
var client *mongo.Client

// Туған күндер сақталған коллекшн
var birthdays *mongo.Collection

func Connect() {
	//URIды файлдан алу
	URI := mongoURI()

	client = connectDB(URI)
	birthdays = client.Database("AsklyBot").Collection("birthdays")
}

func mongoURI() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("can't load .env file")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		log.Fatalf("can't find URI file")
	}

	return MONGO_URI
}

func connectDB(MONGO_URI string) *mongo.Client {
	// URIды жүктеу
	clientOptions := options.Client().ApplyURI(MONGO_URI)

	// Монгодбға қосылу
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Тексеру
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully!")

	// Клиентті қайтарамыз
	return client
}
