package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkerHistory struct {
	Worker      primitive.ObjectID `json:"worker" bson:"worker"`
	Event       string             `json:"event" bson:"event"`
	Description string             `json:"description" bson:"description"`
	Value       bson.M             `json:"value" bson:"value"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
}
