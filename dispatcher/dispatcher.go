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

	IsRunning bool

	Delay int
}

func (d *Dispatcher) Run() {

	d.IsRunning = true

	if d.Ctx == nil {
		d.Ctx = context.Background()
	}

	count := 0

	fmt.Println("[Dispatcher] START RUNNING...")

	for d.IsRunning {
		count++
		fmt.Println("[TIME]: ", count)
		d.ProcessOne(d.Ctx)
		fmt.Println("[DONE]: ", count)

		delayTime := d.Delay

		if delayTime == 0 {
			delayTime = DELAY_TIME
		}

		fmt.Printf("[DELAY]: %vs\n", delayTime)
		time.Sleep(time.Second * time.Duration(delayTime))
	}

	fmt.Println("[Dispatcher] STOP RUNNING...")
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

	err = pusher.CommitData(ctx, data)

	if err != nil {
		log.Println(err)

		return
	}
}
