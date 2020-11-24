package model

import (
	"context"
	"reflect"
	"testing"

	"github.com/barrydevp/codeatest-runner-core/connections"
	"go.mongodb.org/mongo-driver/bson"
)

func TestHeartbeat(t *testing.T) {
	HeartbeatColl := connections.GetModel("heartbeats")

	var heartbeat Heartbeat

	err := HeartbeatColl.FindOne(context.TODO(), bson.M{}).Decode(&heartbeat)

	if err != nil {
		t.Log(reflect.TypeOf(err))
		t.Error(err)
	}

	t.Log(heartbeat)
}
