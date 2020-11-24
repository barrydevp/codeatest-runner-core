package puller

import (
	// "github.com/barrydevp/codeatest-runner-core/connections"
	"github.com/barrydevp/codeatest-runner-core/model"
)

type Submit struct {
	Bucket []model.Submit
	Limit  int64
}

func (submit *Submit) PullSubmit(query interface{}) model.Submit {

	return model.Submit{}
}
