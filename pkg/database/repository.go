package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// Туған күндерді сақтауға арналған формат
type Birthday struct {
	ID       string   `bson:"id"`
	Name     string   `bson:"name"`
	Birthday string   `bson:"birthday"`
	GroupIDs []string `bson:"group_ids"`
}

// Туған күнді қосу үшін
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

// Егер адам базада бар болса, тек жаңа группаға қосады.
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
