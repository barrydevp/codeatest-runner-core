package server

import (
	"encoding/json"
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
	http.HandleFunc("/state", hs.state)

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

func (hs *HttpServer) state(w http.ResponseWriter, r *http.Request) {

	state := map[string]interface{}{
		"dispatcher": map[string]interface{}{
			"name":       hs.Dispatcher.Name,
			"is_running": hs.Dispatcher.IsRunning,
			"delay":      hs.Dispatcher.Delay,
			"run_count":  hs.Dispatcher.RunCount,
			"cycle":      hs.Dispatcher.Cycle,
		},
		"runner": map[string]interface{}{
			"name":    hs.Dispatcher.Runner.Name,
			"state":   hs.Dispatcher.Runner.State,
			"command": hs.Dispatcher.Runner.Command,
		},
		"puller": map[string]interface{}{
			"language":    hs.Dispatcher.Puller.Language,
			"bucket_size": hs.Dispatcher.Puller.BucketSize,
		},
	}
	stateJson, err := json.Marshal(state)

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "application/json")

		w.Write(stateJson)
	}
}
