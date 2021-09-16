package eats

import (
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

var WorkflowActivityOptions = workflow.ActivityOptions{
	TaskList:               "orderTasks",
	ScheduleToCloseTimeout: time.Second * 80,
	ScheduleToStartTimeout: time.Second * 80,
	StartToCloseTimeout:    time.Second * 80,
	HeartbeatTimeout:       time.Second * 20,
	WaitForCancellation:    false,
}
const OrderReceivedSignalToken = "OrderReceivedSignalToken"

func OrderWorkflow(ctx workflow.Context, totalValue string) error {
	log.Println("Workflow started in:", totalValue)
	ctx = workflow.WithActivityOptions(ctx, WorkflowActivityOptions)
	workflow.GetLogger(ctx).Info("OrderWorkflow started", zap.String("order total value:", totalValue))

	// Payment Activity
	totalValueConverted, _ := strconv.Atoi(totalValue)
	var paymentResult map[string]interface{}
	paymentFuture := workflow.ExecuteActivity(ctx, PaymentActivity, totalValueConverted)
	if err := paymentFuture.Get(ctx, &paymentResult); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Payment Result", zap.Any("response:", paymentResult))

	// Save Payment Result
	info := workflow.GetInfo(ctx)
	var persistResult bool
	persistFuture := workflow.ExecuteActivity(ctx, PersistActivity, paymentResult, info.WorkflowExecution.ID)
	if err := persistFuture.Get(ctx, &persistResult); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("PersistActivity Result", zap.Any("persistResult:", persistResult))

	err := workflow.Sleep(ctx, time.Second * 3)
	if err != nil {
		return err
	}

	//Activity ETA + Courier
	var etaResult map[string]interface{}
	etaFuture := workflow.ExecuteActivity(ctx, ETAActivity)
	if err := etaFuture.Get(ctx, &etaResult); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("ETA and Courier Name Result", zap.Any("response:", etaFuture))

	// Save ETACourier Result
	var ETAResult bool
	ETAFuture := workflow.ExecuteActivity(ctx, PersistActivity, etaResult, info.WorkflowExecution.ID)
	if err := ETAFuture.Get(ctx, &ETAResult); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Payment Result", zap.Any("response:", paymentResult))

	// Order Received
	selector := workflow.NewSelector(ctx)
	signalChan := workflow.GetSignalChannel(ctx, OrderReceivedSignalToken)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, nil)
		workflow.GetLogger(ctx).Info("Order received")
	})
	selector.Select(ctx)

	return nil
}
