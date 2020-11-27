package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/barrydevp/codeatest-runner-core/services"
)

func main() {

	fmt.Println("GO --->")

	file := "code/5fac0e502f823e38f0f9c9ac/5fac0c08e113d8633c01edb6/go/3025015a85d4b3de2088c9780ea085b0"

	filePath, err := services.DownloadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	s, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(s))
}
