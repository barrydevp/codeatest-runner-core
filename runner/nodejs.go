package runner

import (
	// "context"
	// "fmt"
	// "os/exec"
	"path/filepath"
	// "time"

	"github.com/barrydevp/codeatest-runner-core/puller"
	"github.com/barrydevp/codeatest-runner-core/model"
)

type Nodejs struct {
    BaseRunner
}

func (g *Nodejs) LoadData(data *puller.Data) error {
    if data == nil {
        return NewError(NilData)
    }

    g.Data = data

    return nil
}

func NewNodejs() Runner {

	abs, err := filepath.Abs("./tests/hello.js")

    if err != nil {
        return nil
    }

    runner := &Nodejs{
        BaseRunner{
            "Nodejs",
            "created",
            &puller.Data{
                model.Job{},
                model.Quiz{},
                model.TestCase{},
                abs,
            },
            nil,
            nil,
            nil,
            []string{"node"},
        },
    }

    return runner
}
