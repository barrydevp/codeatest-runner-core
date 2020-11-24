package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	Name        string             `json:"name"`
	Category    primitive.ObjectID `json:"category"`
	Author      string             `json:"author"`
	Description string             `json:"description"`
	Tags        []string           `json:"tags"`
	Status      string             `json:"status"`
	Password    string             `json:"password"`
	IsDeleted   bool               `json:"is_deleted"`
	StartAt     primitive.DateTime `json:"start_at"`
	ExpiredAt   primitive.DateTime `json:"expired_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at"`
	CreatedAt   primitive.DateTime `json:"created_at"`
}
