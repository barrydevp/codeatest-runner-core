package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserQuiz struct {
	UserTopic primitive.ObjectID `json:"user_topic" bson:"user_topic"`
	Quiz      primitive.ObjectID `json:"quiz" bson:"quiz"`
	Result    Result             `json:"result" bson:"result"`
	CreatedAt primitive.DateTime `json:"created_at"`
}
