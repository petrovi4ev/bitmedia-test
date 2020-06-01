package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	BirthDate string             `bson:"birth_date" json:"birth_date"`
	City      string             `bson:"city" json:"city"`
	Country   string             `bson:"country" json:"country"`
	Email     string             `bson:"email" json:"email"`
	Gender    string             `bson:"gender" json:"gender"`
	LastName  string             `bson:"last_name" json:"last_name"`
}
