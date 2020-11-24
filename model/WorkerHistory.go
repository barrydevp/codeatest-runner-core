package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkerHistory struct {
	Worker      primitive.ObjectID `json:"worker" bson:"worker"`
	Event       string             `json:"event" bson:"event"`
	Description string             `json:"description" bson:"description"`
	Value       interface{}        `json:"event" bson:"event"`
	CreatedAt   primitive.DateTime `json:"created_at"`
}
