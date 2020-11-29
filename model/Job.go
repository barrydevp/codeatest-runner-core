package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobResult struct {
	TestInput  string `json:"test_input" bson:"test_input"`
	TestOutput string `json:"test_output" bson:"test_output"`

	RunOutput   string `json:"run_output" bson:"run_output"`
	RunExitCode int    `json:"run_exit_code" bson:"run_exit_code"`
	RunTime     int64  `json:"run_time" bson:"run_time"`
	RunMemory   int64  `json:"run_memory" bson:"run_memory"` // in kb

	IsPassed bool `json:"is_passed" bson:"is_passed"`
}

type Results []JobResult

type Job struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Worker    primitive.ObjectID `json:"name" bson:"name"`
	Submit    primitive.ObjectID `json:"submit" bson:"submit"`
	Status    string             `json:"status" bson:"status"`
	Results   Results            `json:"results" bson:"results"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}
