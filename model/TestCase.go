package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TestCase struct {
	Quiz   primitive.ObjectID `json:"quiz" bson:"user_quiz"`
	Input  string             `json:"input" bson:"input"`
	Output string             `json:"output" bson:"output"`
}
