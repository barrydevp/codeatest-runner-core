package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Worker struct {
	Name      string             `json:"name" bson:"name"`
	Yaml      string             `json:"yaml" bson:"yaml"`
	Status    string             `json:"status" bson:"status"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
	CreatedAt primitive.DateTime `json:"created_at"`
}
