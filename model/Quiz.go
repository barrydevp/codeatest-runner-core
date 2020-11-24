package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Stat struct {
	Difficult string `json:"difficult" bson:"difficult"`
	Author    string `json:"author" bson:"author"`
	Submitted string `json:"submitted" bson:"submitted"`
	Score     string `json:"score" bson:"score"`
}

type Template struct {
	Lang    string `json:"lang" bson:"lang"`
	Content string `json:"author" bson:"content"`
}

type Limit struct {
	Submit  string `json:"submit" bson:"submit"`
	Memory  string `json:"memory" bson:"memory"`
	Cpu     string `json:"cpu" bson:"cpu"`
	Timeout string `json:"timeout" bson:"timeout"`
}

type Quiz struct {
	Topic       primitive.ObjectID `json:"topic" bson:"topic"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Content     string             `json:"content" bson:"content"`
	Status      string             `json:"status" bson:"status"`
	Templates   []Template         `json:"templates" bson:"templates"`
	Stat        Stat               `json:"stat" bson:"stat"`
	Limit       Limit              `json:"limit" bson:"limit"`
	IsDeleted   bool               `json:"is_deleted" bson:"is_deleted"`
	UpdatedAt   primitive.DateTime `json:"updated_at"`
	CreatedAt   primitive.DateTime `json:"created_at"`
}
