package main

import (
	// "context"
	"fmt"
	// "os/exec"
	// "path/filepath"
	// "time"

	// "github.com/barrydevp/codeatest-runner-core/puller"
	"github.com/barrydevp/codeatest-runner-core/runner"
	// "github.com/barrydevp/codeatest-runner-core/model"
)

func main() {
    run := runner.NewGolang() 

    result, err := run.Run()

    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(result.Input)
        fmt.Println(result.Output)
    }

}
