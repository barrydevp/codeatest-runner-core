package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	abs, err := filepath.Abs("./hello.go")

	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("abs", abs)

	ctx, done := context.WithTimeout(context.Background(), time.Second*1)

	defer done()

	cmd := exec.CommandContext(ctx, "go", "run", abs)

	result, err := cmd.Output()

	if err != nil {
		fmt.Println("RUN COMMAND ERR", err)
	}

	fmt.Println(string(result))

}
