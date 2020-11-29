package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	Name        string             `json:"name" bson:"name"`
	Category    primitive.ObjectID `json:"category" bson:"category"`
	Author      string             `json:"author" bson:"author"`
	Description string             `json:"description" bson:"descripton"`
	Tags        []string           `json:"tags" bson:"tags"`
	Status      string             `json:"status" bson:"status"`
	Password    string             `json:"password" bson:"password"`
	IsDeleted   bool               `json:"is_deleted" bson:"is_deleted"`
	StartAt     primitive.DateTime `json:"start_at" bson:"start_at"`
	ExpiredAt   primitive.DateTime `json:"expired_at" bson:"expired_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
}
