package main

import (
	"context"
	"fmt"
	"log"

	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/puller"
	"github.com/barrydevp/codeatest-runner-core/pusher"
)

func main() {
	ctx := context.Background()

	data, err := puller.PullData(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data.FilePath)

	err = pusher.MarkProcessingSubmit(ctx, data.Submit)

	if err != nil {
		log.Fatal(err)
	}

	data.Submit.Status = "completed"
	data.Submit.Result = model.Result{
		100,
		1,
		1000,
	}

	err = pusher.CommitData(ctx, data)

	if err != nil {
		log.Fatal(err)
	}
}
