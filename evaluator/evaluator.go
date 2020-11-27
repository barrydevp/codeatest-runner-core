package evaluator

import (
	"fmt"
	"syscall"

	"github.com/barrydevp/codeatest-runner-core/runner"
)

type Result struct {
	TestInput  string
	TestOutput string

	RunOutput   string
	RunExitCode int
	RunTime     int64
	RunMemory   int64 // in kb
}

func Evaluate(rCmd *runner.RunnerCmd) bool {
	fmt.Println(rCmd.Cmd.String())

	fmt.Printf("OUT: %sTEST: %s", string(rCmd.Output), rCmd.TestCase.Output)

	usage, _ := rCmd.Cmd.ProcessState.SysUsage().(*syscall.Rusage)

	memory := usage.Maxrss
	time := usage.Utime.Usec

	fmt.Printf("STATS: mem: %v cpu: %v\n", memory, time)

	fmt.Println("Pass: ", rCmd.Output == rCmd.TestCase.Output)

	return true
}
