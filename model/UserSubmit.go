package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserQuiz struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserTopic primitive.ObjectID `json:"user_topic" bson:"user_topic"`
	Quiz      primitive.ObjectID `json:"quiz" bson:"quiz"`
	Result    Result             `json:"result" bson:"result"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`

	UserTopicObj *UserTopic
	QuizObj      *Quiz
}
