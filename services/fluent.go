package services

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

var FluentClient *fluent.Fluent
var FLUENT_HOST string
var FLUENT_PORT int

func init() {

	FLUENT_HOST = os.Getenv("FLUENT_HOST")
	FLUENT_PORT, err := strconv.Atoi(os.Getenv("FLUENT_PORT"))

	if err != nil {
		log.Fatal("Invalid FLUENT_PORT: ", os.Getenv("FLUENT_PORT"))
	}

	if FLUENT_HOST == "" {
		log.Fatal("Missing env: FLUENT_HOST")
	}

	FluentClient, err = fluent.New(fluent.Config{
		FluentHost:    FLUENT_HOST,
		FluentPort:    FLUENT_PORT,
		MarshalAsJSON: true,
		Timeout:       time.Duration(time.Second * 5),
	})

	if err != nil {
		log.Fatal("Cannot connect to fluentd, error: ", err.Error())
	}
}
