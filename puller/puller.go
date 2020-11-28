package puller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Job model.Job

	Quiz model.Quiz

	TestCases []model.TestCase

	FilePath string
}

func PullData() (Data, error) {

	return Data{}, nil
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

func GetSubmit(ctx context.Context) (*model.Submit, error) {
	SubmitColl := connections.GetModel("submits")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))

	defer cancel()

	var submit model.Submit

	opts := options.FindOne().SetSort(bson.D{{"created_at", 1}})

	err := SubmitColl.FindOne(ctxTimeout, bson.M{
		"status": bson.M{"$in": []string{"pending", "retry"}},
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

	quiz, err := GetQuizV2(ctx, userQuiz.Quiz)

	if err != nil {
		return nil, err
	}

	userQuiz.QuizObj = quiz

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

	stage := bson.D{
		{"$match", bson.M{"_id": quizId}},
		{"$lookup", bson.M{
			"from":         "testcases",
			"localField":   "_id",
			"foreignField": "quiz",
			"as":           "test_case_objs",
		}},
	}

	if obj, err := json.Marshal(stage); err == nil {
		fmt.Println(string(obj))
	}

	opts := options.Aggregate().SetMaxTime(2 * time.Second)

	cursor, err := QuizColl.Aggregate(ctx, stage, opts)

	if err != nil {
		return nil, err
	}

	fmt.Println("OK")
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
