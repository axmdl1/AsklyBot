package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID       int       `json:"id" bson:"id"` // Числовой ID
	Name     string    `json:"name" bson:"name"`
	Birthday time.Time `json:"birthday" bson:"birthday"`
}

type RequestJSON struct {
	Message string `json:"message"`
}

type ResponseJSON struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var collection *mongo.Collection

// Connect to MongoDB
func connectToDatabase() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collection = client.Database("appDB").Collection("people")
	fmt.Println("Successfully connected to MongoDB")
}

// Create a new person
func createPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Name     string `json:"name"`
		Birthday string `json:"birthday"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	parsedBirthday, err := time.Parse("2006-01-02", input.Birthday)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Найти максимальный существующий ID
	var lastPerson Person
	err = collection.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"id": -1})).Decode(&lastPerson)
	if err != nil && err != mongo.ErrNoDocuments {
		http.Error(w, "Failed to retrieve last person", http.StatusInternalServerError)
		return
	}

	// Увеличиваем ID на 1 для нового пользователя
	newID := lastPerson.ID + 1

	person := Person{
		ID:       newID,
		Name:     input.Name,
		Birthday: parsedBirthday,
	}

	_, err = collection.InsertOne(ctx, person)
	if err != nil {
		http.Error(w, "Failed to save person", http.StatusInternalServerError)
		return
	}

	// Возвращаем ID созданного пользователя
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseJSON{
		Status:  "success",
		Message: fmt.Sprintf("Person created successfully with ID: %d", newID),
	})
}

// Read all people
func readPeopleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var people []Person
	if err := cursor.All(ctx, &people); err != nil {
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// Read person by ID
func readPersonByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseJSON{
			Status:  "fail",
			Message: "ID parameter is required",
		})
		return
	}

	// Конвертируем ID в int
	intID, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseJSON{
			Status:  "fail",
			Message: "Invalid ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var person Person
	err = collection.FindOne(ctx, bson.M{"id": intID}).Decode(&person)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseJSON{
			Status:  "fail",
			Message: "Person not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func postHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]interface{} // Используем map для проверки наличия ключа
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request["message"] == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseJSON{
			Status:  "fail",
			Message: "Invalid JSON message",
		})
		return
	}

	// Проверяем, что значение "message" — это строка
	message, ok := request["message"].(string)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseJSON{
			Status:  "fail",
			Message: "Invalid JSON message",
		})
		return
	}

	// Если всё корректно, возвращаем успех с самим сообщением
	w.Header().Set("Content-Type", "application/json")
	if message == "" {
		fmt.Println("Message received: (empty)")
		json.NewEncoder(w).Encode(ResponseJSON{
			Status:  "success",
			Message: "Data successfully received (empty message)",
		})
	} else {
		fmt.Println("Message received:", message)
		json.NewEncoder(w).Encode(ResponseJSON{
			Status:  "success",
			Message: fmt.Sprintf("%s", message),
		})
	}
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	// For GET requests, just return an informational message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseJSON{
		Status:  "success",
		Message: "Use POST Method to send data",
	})
}

// Main function
func main() {
	connectToDatabase()

	http.HandleFunc("/create", createPersonHandler)
	http.HandleFunc("/read", readPeopleHandler)
	http.HandleFunc("/readByID", readPersonByIDHandler)
	http.HandleFunc("/post", postHandle)
	http.HandleFunc("/get", getHandle)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
