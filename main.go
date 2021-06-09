package main

import (
	"github.com/menezee/cadence-project/eats"
	"github.com/menezee/cadence-project/helpers"
	"github.com/uber-go/tally"
	_ "go.uber.org/cadence/.gen/go/cadence"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	_ "go.uber.org/yarpc/api/transport"
	"go.uber.org/zap"
)

func main() {
	cadenceClient := helpers.BuildCadenceClient()
	startWorker(helpers.BuildLogger(), cadenceClient)
	select {}
}

func startWorker(logger *zap.Logger, service workflowserviceclient.Interface) {
	workerOptions := worker.Options{
		Logger:       logger,
		MetricsScope: tally.NewTestScope(helpers.TaskListName, map[string]string{}),
	}

	workerInstance := worker.New(
		service,
		helpers.Domain,
		helpers.TaskListName,
		workerOptions)

	workerInstance.RegisterActivity(eats.DebitActivity)
	workerInstance.RegisterActivity(eats.CreditActivity)
	workerInstance.RegisterWorkflow(eats.CourierTipWorkflow)

	err := workerInstance.Start()
	if err != nil {
		panic("Failed to start workerInstance")
	}

	logger.Info("Started Worker.", zap.String("TaskListName:", helpers.TaskListName))
}
