package pusher

import (
	"context"
	"errors"

	"github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MarkProcessingSubmit(ctx context.Context, submit *model.Submit) error {

	SubmitColl := connections.GetModel("submits")

	filter := bson.D{{"_id", submit.Id}}
	update := bson.D{{"$set", bson.M{"status": "processing"}}}

	submit.Status = "processing"

	result, err := SubmitColl.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("cannot found submit to mark processing")
	}

	return nil
}

func CreateJob(ctx context.Context, job *model.Job) (err error) {

	JobColl := connections.GetModel("jobs")

	jsonBsonM, err := ToBsonM(job)

	if err != nil {
		return
	}

	result, err := JobColl.InsertOne(context.TODO(), jsonBsonM)

	if err != nil {
		return
	}

	job.Id = result.InsertedID.(primitive.ObjectID)

	return nil
}

func ToBsonM(val interface{}) (bsonM *bson.M, err error) {
	bsonBytes, err := bson.Marshal(val)
	if err != nil {
		return
	}

	bsonM = new(bson.M)

	err = bson.Unmarshal(bsonBytes, bsonM)

	if err != nil {
		return nil, err
	}

	return
}
