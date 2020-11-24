package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type General struct {
	Permission int64 `json:"permission" bson:"permission"`
}

type UserTopic struct {
	User      string             `json:"user" bson:"user"`
	Topic     primitive.ObjectID `json:"topic" bson:"topic"`
	General   General            `json:"general" bson:"general"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`
	CreatedAt primitive.DateTime `json:"created_at"`
}
