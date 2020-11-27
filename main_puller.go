package main

import (
	"context"
	"fmt"
	"log"

	"github.com/barrydevp/codeatest-runner-core/puller"
)

func main() {
	ctx := context.Background()

	submit, err := puller.GetSubmit(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(submit)
	fmt.Println(*submit.UserQuizObj)
	fmt.Println(*submit.UserQuizObj.QuizObj)
}
