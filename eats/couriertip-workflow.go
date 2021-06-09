package eats

import (
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

var DefaultActivityOptions = workflow.ActivityOptions{
		TaskList:               "tdcTasks",
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
const BankConfirmationSignalToken = "BankConfirmationSignalToken"

func CourierTipWorkflow(ctx workflow.Context, tipAmount int) error {
	ctx = workflow.WithActivityOptions(ctx, DefaultActivityOptions)
	workflow.GetLogger(ctx).Info("CourierTipWorkflow started", zap.Int("Tip amount", tipAmount))

	var debitResult bool
	debitFuture := workflow.ExecuteActivity(ctx, DebitActivity, tipAmount)
	if err := debitFuture.Get(ctx, &debitResult); err != nil {
		return err
	}

	err := workflow.Sleep(ctx, time.Second * 20)
	if err != nil {
		return err
	}

	var creditResult bool
	creditFuture := workflow.ExecuteActivity(ctx, CreditActivity, tipAmount)
	if err := creditFuture.Get(ctx, &creditResult); err != nil {
		return err
	}

	selector := workflow.NewSelector(ctx)
	var bankConfirmationMessage string

	signalChan := workflow.GetSignalChannel(ctx, BankConfirmationSignalToken)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &bankConfirmationMessage)
		workflow.GetLogger(ctx).Info("Received bank confirmation message", zap.String("Message", bankConfirmationMessage))
	})
	selector.Select(ctx)

	return nil
}
