package dispatcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/barrydevp/codeatest-runner-core/evaluator"
	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/puller"
	"github.com/barrydevp/codeatest-runner-core/pusher"
	"github.com/barrydevp/codeatest-runner-core/runner"
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
}

func (d *Dispatcher) StopRun() {
	d.IsRunning = false
}

func (d *Dispatcher) Run() {

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
		d.ProcessOne(d.Ctx)
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

func (d *Dispatcher) ProcessMany(ctx context.Context) {

}

func (d *Dispatcher) ProcessOne(ctx context.Context) {
	data, err := d.Puller.PullData(ctx)

	if err != nil {
		log.Println(err)

		return
	}

	err = pusher.MarkProcessing(ctx, data)

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
