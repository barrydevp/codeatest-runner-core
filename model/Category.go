package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Slug        string             `json:"slug"`
	Status      string             `json:"status"`
	IsDeleted   bool               `json:"is_deleted"`
	UpdatedAt   primitive.DateTime `json:"updated_at"`
	CreatedAt   primitive.DateTime `json:"created_at"`
}
