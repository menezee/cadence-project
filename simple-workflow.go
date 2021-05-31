package main

import (
	"go.uber.org/zap"
	"time"

	"go.uber.org/cadence/workflow"
)

func TDCWorkflow(ctx workflow.Context, value string) error {
	ao := workflow.ActivityOptions{
		TaskList:               "tdcTasks",
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.Sleep(ctx, time.Second * 10)
	if err != nil {
		return err
	}

	future := workflow.ExecuteActivity(ctx, SimpleActivity, value)
	if err := future.Get(ctx, &result); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("SimpleActivity done", zap.String("result", result))

	// SIGNAL START
	const signalName = "signalForTDC"
	selector := workflow.NewSelector(ctx)
	var age int

	signalChan := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &age)
		workflow.GetLogger(ctx).Info("Received age from signal!", zap.String("signal", signalName), zap.Int("value", age))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)

	var ageCheckResult string
	ageCheckFuture := workflow.ExecuteActivity(ctx, AgeCheckActivity, age)
	if err := ageCheckFuture.Get(ctx, &ageCheckResult); err != nil {
		return err
	}

	workflow.GetLogger(ctx).Info("AgeCheckActivity done", zap.String("User age message", ageCheckResult))
	// SIGNAL END

	workflow.GetLogger(ctx).Info("Workflow done", zap.String("result", result))
	return nil
}
