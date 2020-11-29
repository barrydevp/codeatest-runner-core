package main

import (
	"context"
	"fmt"
	"log"

	"github.com/barrydevp/codeatest-runner-core/puller"
)

func main() {
	ctx := context.Background()

	submit, err := puller.GetSubmit(ctx, "js")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(submit)
	fmt.Println(*submit.UserQuizObj)
}
