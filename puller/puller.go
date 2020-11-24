package puller

import (
	// "github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/model"
)

func PullData() (Data, error) {

	return Data{}, nil
}

type Data struct {
	Job model.Job

	Quiz model.Quiz

	TestCase model.TestCase

	FilePath string
}
