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
	err := workflow.Sleep(ctx, time.Minute * 10)
	if err != nil {
		return err
	}
	future := workflow.ExecuteActivity(ctx, SimpleActivity, value)
	if err := future.Get(ctx, &result); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Done", zap.String("result", result))

	return nil
}
