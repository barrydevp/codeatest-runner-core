package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Worker struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Yaml      string             `json:"yaml" bson:"yaml"`
	Status    string             `json:"status" bson:"status"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
