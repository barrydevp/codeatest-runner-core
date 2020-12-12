package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/evaluator"
	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/puller"
	"github.com/barrydevp/codeatest-runner-core/pusher"
	"github.com/barrydevp/codeatest-runner-core/runner"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const DELAY_TIME = 5

type Dispatcher struct {
	Name string

	Runner *runner.Runner
	Puller *puller.Puller

	Ctx context.Context

	Delay int

	IsRunning bool
	RunCount  int
	Cycle     int

	Worker *model.Worker
}

func (d *Dispatcher) Init() {
	WorkerColl := connections.GetModel("workers")

	ctxTimeout, cancel := context.WithTimeout(d.Ctx, time.Second*time.Duration(10))

	defer cancel()

	var worker model.Worker

	err := WorkerColl.FindOne(ctxTimeout, bson.M{
		"name": d.Name,
	}).Decode(&worker)

	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR-TERMINATE] Cannot load worker %s from Database, error: %s\n", d.Name, err.Error()))
	}

	log.Println("[INIT-STEP] Start mark retry all processing submit...")
	result, err := d.markRetryAllProcessingSubmit()

	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR-TERMINATE] mark retry for all processing submit fail, error: %s\n", err.Error()))
	}

	log.Println(fmt.Sprintf("[INIT-STEP] Mark retry %v processing submit", result.ModifiedCount))

	log.Println("[INIT-STEP] Start change Worker's status to 'up'")
	d.Worker = &worker
	err = d.updateWorkerStatus("up")

	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR-TERMINATE] update status worker fail, error: %s\n", err.Error()))
	}

	log.Println(fmt.Sprintf("WORKER %s Init Successfully\n\t# STATUS: %s\n\t# Updated: %s", worker.Name, worker.Status, worker.UpdatedAt))
}

func (d *Dispatcher) updateWorkerStatus(status string) error {
	WorkerColl := connections.GetModel("workers")

	updatedAt := time.Now()
	filter := bson.D{{"_id", d.Worker.Id}}
	update := bson.D{{"$set", bson.M{"status": status, "updated_at": updatedAt}}}

	result, err := WorkerColl.UpdateOne(d.Ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("cannot found worker to change status")
	}

	d.Worker.Status = status
	d.Worker.UpdatedAt = updatedAt

	return nil
}

// Mark pending previous processing job
func (d *Dispatcher) markRetryAllProcessingSubmit() (*mongo.UpdateResult, error) {
	SubmitColl := connections.GetModel("submits")

	updatedAt := time.Now()
	filter := bson.D{{"status", "processing"}, {"language", d.Puller.Language}}
	update := bson.D{{"$set", bson.M{"status": "retry", "updated_at": updatedAt}}}

	result, err := SubmitColl.UpdateOne(d.Ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}
func (d *Dispatcher) StopRun() {
	d.IsRunning = false
}

func (d *Dispatcher) Run() {

	if d.Worker == nil {
		log.Fatal("[Dispatcher] Cannot found worker to run.")
	}

	d.IsRunning = true

	if d.Ctx == nil {
		d.Ctx = context.Background()
	}

	d.RunCount = 0

	fmt.Println("[Dispatcher] START RUNNING...")

	for d.IsRunning {
		if d.RunCount > 100000 {
			d.Cycle++
			d.RunCount = 0
		}
		d.RunCount++
		fmt.Println("[TIME]: ", d.RunCount)
		d.ProcessMany()
		fmt.Println("[DONE]: ", d.RunCount)

		delayTime := d.Delay

		if delayTime == 0 {
			delayTime = DELAY_TIME
		}

		fmt.Printf("[DELAY]: %vs\n", delayTime)
		time.Sleep(time.Second * time.Duration(delayTime))
	}

	fmt.Println("[Dispatcher] STOP RUNNING...")
}

func (d *Dispatcher) ProcessMany() {
	datas, err := d.Puller.PullDatas(d.Ctx, d.Worker)

	if err != nil {
		log.Println(err)

		return
	}

	wg := sync.WaitGroup{}

	log.Println(len(datas))

	fork := func(data *puller.Data) {
		d.ProcessOne(data)

		wg.Done()
	}

	for i, _ := range datas {
		wg.Add(1)
		go fork(&datas[i])
	}

	wg.Wait()

	return
}

func (d *Dispatcher) ProcessOne(data *puller.Data) {
	ctx := d.Ctx

	err := pusher.MarkProcessing(ctx, data)

	if err != nil {
		log.Println(err)

		return
	}

	rCmds, err := d.Runner.Process(data)

	if err != nil {
		log.Println(err)

		return
	}

	results := make(model.Results, 0, len(data.Quiz.TestCaseObjs))

	for _, rCmd := range rCmds {

		result := evaluator.Evaluate(rCmd)

		results = append(results, *result)
	}

	data.Job.Results = results
	data.Job.Status = "done"
	data.Submit.Status = "completed"
	data.Submit.Result = *evaluator.CaculateResult(data.Submit, data.Quiz, results)
	log.Println("[RESULT]: ", data.Submit.Result)

	err = pusher.CommitData(ctx, data)

	if err != nil {
		log.Println(err)

		return
	}
}
