package connections

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestConnect(t *testing.T) {

	User := GetModel("users")

	cursor, err := User.Find(context.Background(), bson.M{})

	if err != nil {
		t.Errorf("connect fail %s", err.Error())
	}

	for cursor.Next(context.TODO()) {
		t.Log(cursor.Current)
	}

}
