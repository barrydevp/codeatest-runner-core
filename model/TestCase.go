package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TestCase struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	Quiz   primitive.ObjectID `json:"quiz" bson:"quiz"`
	Input  string             `json:"input" bson:"input"`
	Output string             `json:"output" bson:"output"`
}
