package puller

import (
	"context"
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

	quiz, err := GetQuiz(ctx, userQuiz.Quiz)

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
