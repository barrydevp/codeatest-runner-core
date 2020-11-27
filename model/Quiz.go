package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Stat struct {
	Difficult string `json:"difficult" bson:"difficult"`
	Author    string `json:"author" bson:"author"`
	Submitted int64  `json:"submitted" bson:"submitted"`
	Score     int64  `json:"score" bson:"score"`
}

type Template struct {
	Lang    string `json:"lang" bson:"lang"`
	Content string `json:"author" bson:"content"`
}

type Limit struct {
	Submit  int64 `json:"submit" bson:"submit"`
	Memory  int64 `json:"memory" bson:"memory"`
	Cpu     int64 `json:"cpu" bson:"cpu"`
	Timeout int64 `json:"timeout" bson:"timeout"`
}

type Quiz struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Topic       primitive.ObjectID `json:"topic" bson:"topic"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Content     string             `json:"content" bson:"content"`
	Status      string             `json:"status" bson:"status"`
	Templates   []Template         `json:"templates" bson:"templates"`
	Stat        Stat               `json:"stat" bson:"stat"`
	Limit       Limit              `json:"limit" bson:"limit"`
	IsDeleted   bool               `json:"is_deleted" bson:"is_deleted"`
	UpdatedAt  primitive.DateTime `json:"updated_at" bson:"updated_at"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`

    TestCaseObjs []TestCase
}
