package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/puller"
)

type Runner struct {
	Name  string
	State string

	Command  string
	BaseArgs []string
	Env      []string
	Dir      string
}

type RunnerCmd struct {
	Cmd *exec.Cmd

	TestCase *model.TestCase

	Limit *model.Limit

	Output string
}

func (r *Runner) Process(data *puller.Data) ([]*RunnerCmd, error) {

	if data == nil {
		return nil, errors.New("[RunnerErro]: invalid data.")
	}

	quiz := data.Quiz

	limit := quiz.Limit

	timeout := limit.Timeout
	// memory := limit.Memory

	var timeoutSec int64 = 10

	if timeout > 0 {
		timeoutSec = timeout
	}

	timeoutDur := time.Second * time.Duration(timeoutSec)

	testCases := quiz.TestCaseObjs

	rCmds := make([]*RunnerCmd, 0, len(testCases))

	for i := 0; i < len(testCases); i++ {
		input := testCases[i].Input

		args := make([]string, 0, len(r.BaseArgs)+1)
		args = append(args, r.BaseArgs...)
		args = append(args, data.FilePath)

		ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
		defer cancel()

		cmd := exec.CommandContext(ctx, r.Command, args...)
		cmd.Dir = r.Dir
		cmd.Env = r.Env

		stdin, err := cmd.StdinPipe()

		if err != nil {
			return nil, errors.New(fmt.Sprintf("[RunnerError]: %s", err.Error()))
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, input)
		}()

		rCmd := RunnerCmd{
			cmd,
			&testCases[i],
			&limit,
			"",
		}

		rCmds = append(rCmds, &rCmd)
	}

	var wg sync.WaitGroup

	fork := func(rCmd *RunnerCmd) {
		cmd := rCmd.Cmd

		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("[RunnerLog]: Error while run cmd %s\n", err.Error())
		}

		rCmd.Output = string(output)

		wg.Done()
	}

	wg.Add(len(rCmds))

	for i := 0; i < len(rCmds); i++ {
		go fork(rCmds[i])
	}

	wg.Wait()

	// clear file
	err := os.Remove(data.FilePath)

	if err != nil {
		log.Printf("[WARNING] Cannot rm %s\n", data.FilePath)
	}

	return rCmds, nil
}
