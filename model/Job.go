package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Worker    primitive.ObjectID `json:"name" bson:"name"`
	Submit    primitive.ObjectID `json:"submit" bson:"submit"`
	Status    string             `json:"status" bson:"status"`
	Result    primitive.M        `json:"result" bson:"result"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}
