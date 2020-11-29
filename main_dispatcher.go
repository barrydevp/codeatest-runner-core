package main

import (
	"context"

	"github.com/barrydevp/codeatest-runner-core/dispatcher"
	"github.com/barrydevp/codeatest-runner-core/puller"
	"github.com/barrydevp/codeatest-runner-core/runner"
)

func main() {

	NodeJSRunner := runner.Runner{
		Name:    "NodeJS",
		State:   "created",
		Command: "node",
	}

	NodeJSPuller := puller.Puller{
		Language:   "js",
		BucketSize: 1,
	}

	ctx := context.Background()

	NodeJSDispatcher := dispatcher.Dispatcher{
		Name:      "NodeJS",
		Runner:    &NodeJSRunner,
		Puller:    &NodeJSPuller,
		Ctx:       ctx,
		IsRunning: false,
	}

	NodeJSDispatcher.Run()
}
