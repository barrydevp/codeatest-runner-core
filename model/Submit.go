package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Result struct {
	Score      int32 `json:"score" bson:"score"`
	Time       int64 `json:"time" bson:"time"`
	MemoryUsed int64 `json:"memory_used" bson:"memory_used"`
}

type Submit struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	UserQuiz   primitive.ObjectID `json:"user_quiz" bson:"user_quiz"`
	Language   string             `json:"language" bson:"language"`
	UploadFile string             `json:"upload_file" bson:"upload_file"`
	Status     string             `json:"status" bson:"status"`
	Result     Result             `json:"result" bson:"result"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`

	UserQuizObj *UserQuiz `json:"user_quiz_obj" bson:"user_quiz_obj"`

	UserQuizObjs []UserQuiz `json:"user_quiz_objs" bson:"user_quiz_objs"`
}
