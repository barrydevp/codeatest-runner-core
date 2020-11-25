package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

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
	Run() error

	Reload() error

	Clear() error

	Kill() error

	GetName() string

	GetContext() *context.Context

	GetState() string

	GetData() *puller.Data

	GetCommand() *exec.Cmd

	GetResult() *Result

	GetEvents() []Event
}

type BaseRunner struct {
	Name string

	State string

	Data *puller.Data

	Result *Result

	Ctx *context.Context

	Events []Event

	Command *exec.Cmd
}

const (
	ErrorLabel  = "RunnerError"
	FormatError = "[%s]: %s"

	NilData = "nil data."

	NilCommand = "nil command."
)

func NewError(message string) error {

	return errors.New(fmt.Sprintf(FormatError, ErrorLabel, message))
}

func (this *BaseRunner) LoadData(data *puller.Data) error {
	if data == nil {
        return NewError(NilData)
	}

	this.Data = data

    this.Result = &Result{
        data.TestCase.Input,
        "",
    }

	return nil
}

func (this *BaseRunner) Run() error {
	if this.Command == nil {
        return NewError(NilCommand)
	}

    output, err := this.Command.Output()

    if err != nil {
        return NewError(err.Error())
    }

    this.Result.Output = string(output)

	return nil
}

func (this *BaseRunner) GetName() string {
	return this.Name
}

func (this *BaseRunner) GetState() string {
	return this.State
}

func (this *BaseRunner) GetData() *puller.Data {
	return this.Data
}

func (this *BaseRunner) GetResult() *Result {
	return this.Result
}

func (this *BaseRunner) GetContext() *context.Context {
	return this.Ctx
}

func (this *BaseRunner) GetEvents() []Event {
	return this.Events
}

func (this *BaseRunner) GetCommand() *exec.Cmd {
	return this.Command
}
