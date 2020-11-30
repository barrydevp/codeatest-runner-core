package pusher

import (
	"context"
	"errors"

	"github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/puller"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MarkProcessing(ctx context.Context, data *puller.Data) (err error) {
	err = MarkProcessingSubmit(ctx, data.Submit)

	if err != nil {
		return
	}

	err = CreateJob(ctx, data.Job)

	if err != nil {
		return
	}

	return nil
}

func MarkProcessingSubmit(ctx context.Context, submit *model.Submit) error {

	SubmitColl := connections.GetModel("submits")

	filter := bson.D{{"_id", submit.Id}}
	update := bson.D{{"$set", bson.M{"status": "processing"}}}

	result, err := SubmitColl.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("cannot found submit to mark processing")
	}

	submit.Status = "processing"

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

func CommitData(ctx context.Context, data *puller.Data) (err error) {

	JobColl := connections.GetModel("jobs")
	SubmitColl := connections.GetModel("submits")

	job := data.Job
	submit := data.Submit

	filterSubmit := bson.D{{"_id", submit.Id}}
	updateSubmit := bson.D{{"$set", bson.M{"status": submit.Status, "result": submit.Result}}}

	result, err := SubmitColl.UpdateOne(ctx, filterSubmit, updateSubmit)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("cannot commit for submit")
	}

	filterJob := bson.D{{"_id", job.Id}}
	updateJob := bson.D{{"$set", bson.M{"status": job.Status, "results": job.Results}}}

	result, err = JobColl.UpdateOne(ctx, filterJob, updateJob)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("cannot commit for job")
	}

	err = updateResultUserQuiz(data.Submit, data.Submit.UserQuizObj)

	if err != nil {
		return err
	}

	return nil
}

func updateResultUserQuiz(submit *model.Submit, userQuiz *model.UserQuiz) error {

	UserQuizColl := connections.GetModel("userquizzes")

	lastResult := userQuiz.Result
	newResult := submit.Result

	if lastResult.Score > newResult.Score {
		return nil
	}

	filter := bson.D{{"_id", submit.UserQuiz}}
	update := bson.D{{"$set", bson.M{"result": bson.M{"score": newResult.Score, "memory_used": newResult.MemoryUsed, "time": newResult.Time}}}}

	result, err := UserQuizColl.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("cannot found user_quiz to mark processing")
	}

	submit.Status = "processing"

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

func ToBsonD(val interface{}) (bsonD *bson.D, err error) {
	bsonBytes, err := bson.Marshal(val)
	if err != nil {
		return
	}

	bsonD = new(bson.D)

	err = bson.Unmarshal(bsonBytes, bsonD)

	if err != nil {
		return nil, err
	}

	return
}
