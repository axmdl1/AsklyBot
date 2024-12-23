package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	//"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var birthdays *mongo.Collection

type Birthday struct {
	ID       string    `bson:"id"`
	Name     string    `bson:"name"`
	Birthday time.Time `bson:"birthday"`
	GroupIDs []string  `bson:"group_ids"`
}

func Connect() {
	URI := mongoURI()

	client = connectDB(URI)
	birthdays = client.Database("AsklyBot").Collection("birthdays")
}

func AddBirthday(birthday Birthday) error {
	var existingUser Birthday
	err := birthdays.FindOne(context.Background(), bson.M{"id": birthday.ID}).Decode(&existingUser)
	if err == nil {
		return addGroup(birthday.ID, birthday.GroupIDs)
	}

	_, err = birthdays.InsertOne(context.Background(), birthday)
	if err != nil {
		log.Fatalf("Error inserting document: %v", err)
		return err
	}

	return nil
}

func addGroup(userID string, groupID []string) error {
	_, err := birthdays.UpdateOne(
		context.Background(),
		bson.M{"id": userID},
		bson.M{"$addToSet": bson.M{"group_ids": bson.M{"$each": groupID}}},
	)
	if err != nil {
		log.Fatalf("Error updating user: %v", err)
		return err
	}

	return nil
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
	// Настройка параметров подключения
	clientOptions := options.Client().ApplyURI(MONGO_URI)

	// Попытка подключения к MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Проверка на успешное подключение (Ping)
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully!")

	// Возврат подключенного клиента
	return client
}
