package evaluator

import (
	"fmt"
	"syscall"

	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/runner"
)

func Evaluate(rCmd *runner.RunnerCmd) *model.JobResult {
	// fmt.Println(rCmd.Cmd.String())

	testOutput := rCmd.TestCase.Output
	testInput := rCmd.TestCase.Input
	runOutput := ""
	if rCmd.Output != "" {
		runOutput = string(rCmd.Output[:len(rCmd.Output)-1])
	}

	fmt.Printf("TEST: -- IN: %s \nOUT: %s\n", testInput, testOutput)
	fmt.Printf("RUN: OUT: %s\n", runOutput)

	var memory int64 = 0
	var time int64 = 0

	exitCode := rCmd.Cmd.ProcessState.ExitCode()

	if exitCode == 0 {
		usage, _ := rCmd.Cmd.ProcessState.SysUsage().(*syscall.Rusage)

		memory = usage.Maxrss
		time = usage.Utime.Usec
	}

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

func CaculateResult(submit *model.Submit, quiz *model.Quiz, results []model.JobResult) *model.Result {
	lastResult := submit.Result
	quizScore := quiz.Stat.Score

	numberOfTest := len(quiz.TestCaseObjs)
	if numberOfTest == 0 {
		numberOfTest = 1
	}

	scorePerTest := quizScore / int32(numberOfTest)

	var totalScore int32 = 0
	var totalMemoryUsed int64 = 0
	var totalTime int64 = 0

	for index, _ := range results {
		totalMemoryUsed += results[index].RunMemory
		totalTime += results[index].RunTime

		if results[index].IsPassed {
			totalScore += scorePerTest
		}
	}

	result := model.Result{
		lastResult.Score,
		lastResult.Time,
		lastResult.MemoryUsed,
	}

	if lastResult.Score > totalScore {
		return &result
	}

	if lastResult.Score == totalScore {
		if lastResult.Time < totalTime {
			return &result
		}

		if lastResult.Time == totalTime {
			if lastResult.MemoryUsed < totalMemoryUsed {
				return &result
			}
		}
	}

	result.Score = totalScore
	result.MemoryUsed = totalMemoryUsed
	result.Time = totalTime

	return &result
}
