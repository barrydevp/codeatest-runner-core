package evaluator

import (
	"fmt"
	"syscall"

	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/runner"
)

func Evaluate(rCmd *runner.RunnerCmd) *model.JobResult {
	fmt.Println(rCmd.Cmd.String())

	testOutput := rCmd.TestCase.Output
	testInput := rCmd.TestCase.Input
	runOutput := string(rCmd.Output[:len(rCmd.Output)-1])

	fmt.Printf("TEST: -- IN: %s \nOUT: %s\n", testInput, testOutput)
	fmt.Printf("RUN: OUT: %s\n", runOutput)

	usage, _ := rCmd.Cmd.ProcessState.SysUsage().(*syscall.Rusage)
	exitCode := rCmd.Cmd.ProcessState.ExitCode()

	memory := usage.Maxrss
	time := usage.Utime.Usec

	fmt.Printf("STATS: mem: %v cpu: %v\n", memory, time)

	isPassed := testOutput == runOutput

	fmt.Println("PASSED: ", isPassed)

	return &model.JobResult{
		testInput,
		testOutput,
		runOutput,
		exitCode,
		time,
		memory,
		isPassed,
	}
}
