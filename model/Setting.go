package model

import "go.mongodb.org/mongo-driver/bson"

type Setting struct {
	Key       string `json:"key" bson"key"`
	Value     bson.M `json:"value" bson:"value"`
	IsDeleted bool   `json:"is_deleted" bson:"is_deleted"`
}
