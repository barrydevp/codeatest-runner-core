package main

import (
	// "context"
	"fmt"
	// "os/exec"
	"path/filepath"
	// "time"

	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/puller"
	"github.com/barrydevp/codeatest-runner-core/runner"
)

func runTest(runner *runner.Runner, filePath string) {
	rCmds, err := runner.Process(
		&puller.Data{
			model.Job{},
			model.Quiz{},
			[]model.TestCase{
				model.TestCase{
					Input:  "",
					Output: "Hello World!\n",
				},
				model.TestCase{
					Input:  "",
					Output: "Hello Codeatest!\n",
				},
			},
			filePath,
		},
	)

	if err != nil {
		fmt.Println(err)
	} else {
		for _, rCmd := range rCmds {
			fmt.Println(rCmd.Cmd.String())

			fmt.Printf("OUT: %sTEST: %s", string(rCmd.Output), rCmd.TestCase.Output)

			fmt.Println("Pass: ", rCmd.Output == rCmd.TestCase.Output)
		}
	}
}

func testNodeJS() {
	run := runner.Runner{
		"NodeJS",
		"created",
		"node",
		[]string{},
	}

	abs, err := filepath.Abs("./tests/hello.js")

	if err != nil {
		fmt.Println(err)

		return
	}

	runTest(&run, abs)
}

func testGolang() {
	run := runner.Runner{
		"Golang",
		"created",
		"go",
		[]string{"run"},
	}

	abs, err := filepath.Abs("./tests/hello.go")

	if err != nil {
		fmt.Println(err)

		return
	}

	runTest(&run, abs)
}

func testPython() {
	run := runner.Runner{
		"Python",
		"created",
		"python3.8",
		[]string{},
	}

	abs, err := filepath.Abs("./tests/hello.py")

	if err != nil {
		fmt.Println(err)

		return
	}

	runTest(&run, abs)
}

func main() {
	testPython()
	testNodeJS()

}
