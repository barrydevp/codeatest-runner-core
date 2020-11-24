package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Heartbeat struct {
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}
