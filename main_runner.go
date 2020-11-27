package main

import (
	// "context"
	"fmt"
	// "syscall"

	// "os/exec"
	"path/filepath"
	// "time"

	"github.com/barrydevp/codeatest-runner-core/evaluator"
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
					Input:  "Hello World!",
					Output: "Hello World!\n",
				},
				model.TestCase{
					Input:  "Hello Codeatest!",
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
			evaluator.Evaluate(rCmd)
			// fmt.Println(rCmd.Cmd.String())

			// fmt.Printf("OUT: %sTEST: %s", string(rCmd.Output), rCmd.TestCase.Output)

			// fmt.Println("Pass: ", rCmd.Output == rCmd.TestCase.Output)

			// usage, _ := rCmd.Cmd.ProcessState.SysUsage().(*syscall.Rusage)

			// usage.Memory
		}
	}
}

func testNodeJS() {
	run := runner.Runner{
		Name:    "NodeJS",
		State:   "created",
		Command: "node",
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
		Name:     "Golang",
		State:    "created",
		Command:  "go",
		BaseArgs: []string{"run"},
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
		Name:    "Python",
		State:   "created",
		Command: "python3.8",
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
	testGolang()
}
