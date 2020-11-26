package runner

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/barrydevp/codeatest-runner-core/puller"
)

type Result struct {
	Input string

	Output string
}

type Event struct {
	Name string

	Value string
}

type Runner interface {
	// Load Data need to run.
	LoadData(*puller.Data) error

	// RUN RUN RUN...
	Run() (*Result, error)

	// Reload() error

	// Kill() error

	GetName() string

	GetState() string

	// GetData() *puller.Data

	// GetCommand() *exec.Cmd

	GetResult() *Result

	GetEvents() []Event
}

type BaseRunner struct {
	Name string

	State string

	Data *puller.Data

	Result *Result

	Events []Event

	Command *exec.Cmd

    BaseArgs []string
}

const (
	ErrorLabel  = "RunnerError"
	FormatError = "[%s]: %s"

	NilData = "nil data."

	NilCommand = "nil command."

	MissingBaseArgs = "missing command args to run."
)

func NewError(message string) error {

	return errors.New(fmt.Sprintf(FormatError, ErrorLabel, message))
}

func (br *BaseRunner) AbleToRun() error {
    if br.Data == nil { 
        return NewError(NilData)
    }
   
    // if br.Command == nil {
    //     return NewError(NilCommand)
    // }

    if len(br.BaseArgs) <= 0 {
        return NewError(MissingBaseArgs)
    }

    return nil
}

// func (br *BaseRunner) CreateResult() *Result {

//     input := br.Data.Quiz.TestCase

// }

func (br *BaseRunner) Run() (*Result, error) {
    err := br.AbleToRun()

    if err != nil {
        return nil, err
    }

    limit := br.Data.Quiz.Limit

    timeout := limit.Timeout
    // memory := limit.Memory

    var timeoutSec int64 = 10

    if(timeout > 0) {
        timeoutSec = timeout
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second * time.Duration(timeoutSec))

    defer cancel()

	// abs, err := filepath.Abs("./tests/hello.go")

    // if err != nil {
        // return nil, err
    // }
    
    cmdArgs := append(br.BaseArgs, br.Data.FilePath)

    cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)

    br.Command = cmd

    output, err := cmd.Output()

    if err != nil {

        return nil, err
    }

    return &Result{
        "",
        string(output),
    }, nil
}

func (br *BaseRunner) GetName() string {
	return br.Name
}

func (br *BaseRunner) GetState() string {
	return br.State
}

func (br *BaseRunner) GetData() *puller.Data {
	return br.Data
}

func (br *BaseRunner) GetResult() *Result {
	return br.Result
}

func (br *BaseRunner) GetEvents() []Event {
	return br.Events
}

func (br *BaseRunner) GetCommand() *exec.Cmd {
	return br.Command
}
