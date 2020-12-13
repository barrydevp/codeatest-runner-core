package helpers

import (
	"log"
	"time"

	"github.com/barrydevp/codeatest-runner-core/model"
	"github.com/barrydevp/codeatest-runner-core/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoggerData struct {
	WorkerId   primitive.ObjectID     `msg:"worker_id"`
	WorkerName string                 `msg:"worker_name"`
	Event      string                 `msg:"event"`
	Status     string                 `msg:"status"`
	Data       map[string]interface{} `msg:"data"`
	CreatedAt  time.Time              `msg:"created_at"`
}

func FluentPost(data *LoggerData) {
	tag := "forward-runner-log"

	err := services.FluentClient.Post(tag, *data)

	if err != nil {
		log.Println("Cannot Log to fluentd, error: ", err.Error())
	}
}

func LogError(worker *model.Worker, event string, data map[string]interface{}) {
	logData := LoggerData{
		WorkerId:   worker.Id,
		WorkerName: worker.Name,
		Event:      event,
		Status:     "error",
		Data:       data,
		CreatedAt:  time.Now(),
	}

	FluentPost(&logData)
}

func LogInfo(worker *model.Worker, event string, data map[string]interface{}) {
	logData := LoggerData{
		WorkerId:   worker.Id,
		WorkerName: worker.Name,
		Event:      event,
		Status:     "info",
		Data:       data,
		CreatedAt:  time.Now(),
	}

	FluentPost(&logData)
}
