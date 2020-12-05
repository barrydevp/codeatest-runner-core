package puller

import (
	"context"
	"errors"
	// "fmt"
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

	BucketSize int64
}

func (p *Puller) PullDatas(ctx context.Context, worker *model.Worker) ([]Data, error) {

	submits, err := GetSubmits(ctx, p.Language, p.BucketSize)

	if err != nil {
		return nil, err
	}

	quizIds := make([]primitive.ObjectID, 0, len(submits))

	datas := make([]Data, 0, len(submits))

	for i, _ := range submits {
		submits[i].UserQuizObj = &submits[i].UserQuizObjs[0]
		quizIds = append(quizIds, submits[i].UserQuizObj.Quiz)

		job := CreateJob(&submits[i], worker)
		filePath, err := p.GetFilePath(&submits[i])

		if err != nil {
			return nil, err
		}

		datas = append(datas, Data{
			Job:      job,
			Submit:   &submits[i],
			FilePath: filePath,
		})
	}

	quizzes, err := GetQuizzes(ctx, quizIds)

	if err != nil {
		return nil, err
	}

	quizzesMap := make(map[primitive.ObjectID]*model.Quiz)

	for i, _ := range quizzes {
		quizzesMap[quizzes[i].Id] = &quizzes[i]
	}

	for i, _ := range datas {
		datas[i].Quiz = quizzesMap[submits[i].UserQuizObj.Quiz]
	}

	return datas, nil
}

func (p *Puller) PullData(ctx context.Context, worker *model.Worker) (*Data, error) {

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

	job := CreateJob(submit, worker)

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

func CreateJob(submit *model.Submit, worker *model.Worker) *model.Job {

	job := model.Job{
		primitive.NilObjectID,
		worker.Id,
		submit.Id,
		"process",
		model.Results{},
		time.Now(),
		time.Now(),
	}

	return &job
}

func GetSubmits(ctx context.Context, language string, limit int64) ([]model.Submit, error) {
	SubmitColl := connections.GetModel("submits")

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(10))

	defer cancel()

	var submits []model.Submit

	matchS := bson.D{
		{"$match", bson.D{{"status", bson.M{"$in": []string{"pending", "retry"}}}, {"language", language}}},
	}
	limitS := bson.D{
		{"$limit", limit},
	}
	lookupS := bson.D{
		{"$lookup", bson.D{
			{"from", "userquizzes"},
			{"localField", "user_quiz"},
			{"foreignField", "_id"},
			{"as", "user_quiz_objs"},
		}},
	}

	cursor, err := SubmitColl.Aggregate(ctxTimeout, mongo.Pipeline{matchS, limitS, lookupS})

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

	opts := options.Aggregate().SetMaxTime(5 * time.Second)

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

func GetQuizzes(ctx context.Context, quizIds []primitive.ObjectID) ([]model.Quiz, error) {
	QuizColl := connections.GetModel("quizzes")

	stage := mongo.Pipeline{
		bson.D{
			{"$match", bson.M{"_id": bson.M{"$in": quizIds}}},
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

	opts := options.Aggregate().SetMaxTime(10 * time.Second)

	cursor, err := QuizColl.Aggregate(ctx, stage, opts)

	if err != nil {
		return nil, err
	}

	var quizzes []model.Quiz

	if err = cursor.All(ctx, &quizzes); err != nil {
		return nil, err
	}

	return quizzes, nil
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
