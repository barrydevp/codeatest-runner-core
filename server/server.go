package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/dispatcher"
)

var RUNNER_NAME string

func init() {
	RUNNER_NAME = os.Getenv("RUNNER_NAME")

	if RUNNER_NAME == "" {
		log.Fatal("Missing env: RUNNER_NAME")
	}
}

type HttpServer struct {
	Dispatcher *dispatcher.Dispatcher
	PORT       string
}

func (hs *HttpServer) ListenAndServe() {
	http.HandleFunc("/", greet)
	http.HandleFunc("/ping", ping)

	if hs.PORT == "" {
		log.Fatal("Missing Server PORT")
	}

	log.Printf("LISTENING ON PORT: %s\n", hs.PORT)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", hs.PORT), nil))
}

func greet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hi i'm %s - %s", RUNNER_NAME, time.Now())
}

func ping(w http.ResponseWriter, r *http.Request) {
	err := connections.Ping()

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Pong! %s", time.Now())
	}
}
