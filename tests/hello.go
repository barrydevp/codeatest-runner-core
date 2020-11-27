package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// s := make([]byte, 0)
	reader := bufio.NewReader(os.Stdin)
	// fmt.Scanln(&s)
	s, _ := reader.ReadString('\n')
	fmt.Println(string(s))
	// fmt.Print("Hello World!")
}
