package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/puller"
)

type Runner struct {
	Name  string `msg:"name"`
	State string `msg:"state"`

	Command  string   `msg:"command"`
	BaseArgs []string `msg:"base_args"`
	Env      []string `msg:"env"`
	Dir      string   `msg:"dir"`

	NeedBuild     bool     `msg:"need_build"`
	BuildCommand  string   `msg:"build_command"`
	BuildBaseArgs []string `msg:"build_base_args"`
}

type RunnerCmd struct {
	Cmd *exec.Cmd

	TestCase *model.TestCase `msg:"test_case"`

	Limit *model.Limit `msg:"limit"`

	Output string `msg:"output"`
}

func (r *Runner) Process(data *puller.Data) ([]*RunnerCmd, error) {
	r.State = "in-processing"

	if data == nil {
		return nil, errors.New("[RunnerError]: invalid data.")
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

	runPath := data.FilePath

	rCmds := make([]*RunnerCmd, 0, len(testCases))

	if r.NeedBuild {
		args := make([]string, 0, len(r.BuildBaseArgs)+1)
		args = append(args, r.BuildBaseArgs...)

		dirPath := filepath.Dir(data.FilePath)
		fileName := filepath.Base(data.FilePath)
		buildPath := filepath.Join(dirPath, strings.Replace(fileName, ".", "", -1))
		args = append(args, buildPath)
		args = append(args, data.FilePath)
		runPath = buildPath

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(120))
		defer cancel()

		cmd := exec.CommandContext(ctx, r.BuildCommand, args...)
		cmd.Dir = r.Dir
		cmd.Env = r.Env

		// fmt.Println(cmd.String())
		output, err := cmd.CombinedOutput()

		if err != nil || cmd.ProcessState.ExitCode() != 0 {
			log.Printf("[RunnerLog]: Error while run build cmd %s\n", err.Error())

			for i := 0; i < len(testCases); i++ {
				rCmd := RunnerCmd{
					cmd,
					&testCases[i],
					&limit,
					string(output),
				}

				rCmds = append(rCmds, &rCmd)
			}

			return rCmds, nil
		}

	}

	for i := 0; i < len(testCases); i++ {
		input := testCases[i].Input

		args := make([]string, 0, len(r.BaseArgs)+1)
		args = append(args, r.BaseArgs...)

		runCommand := runPath

		if !r.NeedBuild {
			args = append(args, runPath)
			runCommand = r.Command
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
		defer cancel()

		cmd := exec.CommandContext(ctx, runCommand, args...)
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
	RemoveFile(data.FilePath)
	RemoveFile(runPath)

	r.State = "after-processing"

	return rCmds, nil
}

func RemoveFile(filePath string) {
	err := os.Remove(filePath)

	if err != nil {
		log.Printf("[WARNING] Cannot rm %s\n", filePath)
	}
}
