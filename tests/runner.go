package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type Test struct {
	name string
}

func (this *Test) GetName() string {
	return this.name
}

type TestChild struct {
	Test

	age int
}

func (this *TestChild) GetName() string {
	return this.name + "!"
}

func (this *TestChild) GetAge() int {
	return this.age
}

func main() {
	abs, err := filepath.Abs("./hello.go")

	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("abs", abs)

	ctx, done := context.WithTimeout(context.Background(), time.Second*1)

	defer done()

	cmd := exec.CommandContext(ctx, "go", "run", abs)

	// fmt.Println(cmd.ProcessState.SysUsage())

	result, err := cmd.Output()

	usage, ok := cmd.ProcessState.SysUsage().(*syscall.Rusage)

	if ok {
		fmt.Println(usage.Maxrss)
		fmt.Println(usage.Utime)
		fmt.Println(usage.Maxrss)
	}

	if err != nil {
		fmt.Println("RUN COMMAND ERR", err)
	}

	fmt.Println(string(result))

	test := TestChild{
		Test{"Hai"},
		12,
	}

	fmt.Println(test.GetName())

}
