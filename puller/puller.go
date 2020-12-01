package puller

import (
	"context"
	"errors"
	"time"

	"github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Job *model.Job

	Submit *model.Submit

	Quiz *model.Quiz

	FilePath string
}

type Puller struct {
	Language string

	BucketSize int
}

func (p *Puller) PullData(ctx context.Context) (*Data, error) {

	submit, err := GetSubmit(ctx, p.Language)

	if err != nil {
		return nil, err
	}

	if submit.UserQuizObj == nil {
		return nil, errors.New("Nil UserQuiz")
	}

	filePath, err := p.GetFilePath(submit)

	if err != nil {
		return nil, err
	}

	quiz, err := GetQuizV2(ctx, submit.UserQuizObj.Quiz)

	if err != nil {
		return nil, err
	}

	job := CreateJob(submit)

	return &Data{
		job,
		submit,
		quiz,
		filePath,
	}, nil
}

func (p *Puller) GetFilePath(submit *model.Submit) (string, error) {

	return services.DownloadFile(submit.UploadFile)
}

func CreateJob(submit *model.Submit) *model.Job {

	workerId, _ := primitive.ObjectIDFromHex("5f942a3588e2242efc747de5")

	job := model.Job{
		primitive.NilObjectID,
		workerId,
		submit.Id,
		"process",
		model.Results{},
		time.Now(),
		time.Now(),
	}

	return &job
}

// func GetSubmits(ctx context.Context) ([]model.Submit, error) {
// 	stage := bson.D{
// 		{"$match", bson.D{{"status", bson.M{"$in": []string{"pending", "retry"}}}}},
// 		{"$limit", 10},
// 		{"$lookup", bson.D{
// 			{"from", "userquizzes"},
// 			{"let", bson.D{{"user_quiz_id", "$user_quiz"}}},
// 			{"pipeline", bson.D{
// 				{"$match", bson.D{{"_id", "$$user_quiz_id"}}},
// 				{"$lookup", bson.D{
// 					{"from", "usertopics"},
// 					{"let", bson.D{{"user_topic_id", "$user_topic"}}},
// 					{"pipeline", bson.D{
// 						{"$match", bson.D{{"_id", "$$user_topic_id"}}},
// 						{"$lookup", bson.D{
// 							{"from", "quizzes"},
// 							{"localField", "quiz"},
// 							{"foreignField", "_id"},
// 							{"as", "quiz"},
// 						}},
// 					}}, }},
// 			}},
// 		}},
// 	}
// }

func GetSubmits(ctx context.Context, language string, limit int64) ([]model.Submit, error) {
	SubmitColl := connections.GetModel("submits")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))

	defer cancel()

	opts := options.Find().SetSort(bson.D{{"created_at", 1}}).SetLimit(limit)

	var submits []model.Submit

	cursor, err := SubmitColl.Find(ctxTimeout, bson.M{
		"status":   bson.M{"$in": []string{"pending", "retry"}},
		"language": language,
	}, opts)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &submits); err != nil {
		return nil, err
	}

	return submits, nil
}

func GetSubmit(ctx context.Context, language string) (*model.Submit, error) {
	SubmitColl := connections.GetModel("submits")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))

	defer cancel()

	var submit model.Submit

	opts := options.FindOne().SetSort(bson.D{{"created_at", 1}})

	err := SubmitColl.FindOne(ctxTimeout, bson.M{
		"status":   bson.M{"$in": []string{"pending", "retry"}},
		"language": language,
	}, opts).Decode(&submit)

	if err != nil {
		return nil, err
	}

	userQuiz, err := GetUserQuiz(ctx, submit.UserQuiz)

	if err != nil {
		return nil, err
	}

	submit.UserQuizObj = userQuiz

	return &submit, nil
}

func GetUserQuiz(ctx context.Context, _id primitive.ObjectID) (*model.UserQuiz, error) {
	UserQuizColl := connections.GetModel("userquizzes")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))

	defer cancel()

	var userQuiz model.UserQuiz

	opts := options.FindOne().SetSort(bson.D{{"created_at", 1}})

	err := UserQuizColl.FindOne(ctxTimeout, bson.M{
		"_id": _id,
	}, opts).Decode(&userQuiz)

	if err != nil {
		return nil, err
	}

	return &userQuiz, nil
}

func GetQuiz(ctx context.Context, _id primitive.ObjectID) (*model.Quiz, error) {
	QuizColl := connections.GetModel("quizzes")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))

	defer cancel()

	var quiz model.Quiz

	opts := options.FindOne().SetSort(bson.D{{"created_at", 1}})

	err := QuizColl.FindOne(ctxTimeout, bson.M{
		"_id": _id,
	}, opts).Decode(&quiz)

	if err != nil {
		return nil, err
	}

	testCases, err := GetTestCases(ctx, quiz.Id)

	quiz.TestCaseObjs = testCases

	return &quiz, nil
}

func GetQuizV2(ctx context.Context, quizId primitive.ObjectID) (*model.Quiz, error) {
	QuizColl := connections.GetModel("quizzes")

	stage := mongo.Pipeline{
		bson.D{
			{"$match", bson.M{"_id": quizId}},
		},
		bson.D{
			{"$lookup", bson.M{
				"from":         "testcases",
				"localField":   "_id",
				"foreignField": "quiz",
				"as":           "test_case_objs",
			}},
		},
	}

	opts := options.Aggregate().SetMaxTime(2 * time.Second)

	cursor, err := QuizColl.Aggregate(ctx, stage, opts)

	if err != nil {
		return nil, err
	}

	var quizzes []model.Quiz

	if err = cursor.All(ctx, &quizzes); err != nil {
		return nil, err
	}

	return &quizzes[0], nil
}

func GetTestCases(ctx context.Context, quizId primitive.ObjectID) ([]model.TestCase, error) {
	TestCaseColl := connections.GetModel("testcases")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))

	defer cancel()

	var testCases []model.TestCase

	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := TestCaseColl.Find(ctxTimeout, bson.M{
		"quiz": quizId,
	}, opts)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &testCases); err != nil {
		return nil, err
	}

	return testCases, nil
}
