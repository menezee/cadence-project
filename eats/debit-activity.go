package eats

import (
	"context"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func DebitActivity(ctx context.Context, amount int) (bool, error) {
	activity.GetLogger(ctx).Info("DebitActivity started.", zap.Int("Amount", amount))
	accountTransactionService := NewAccountTransactionService()
	ok := accountTransactionService.Debit(amount)

	if ok {
		activity.GetLogger(ctx).Info("DebitActivity completed.", zap.Bool("Succeeded", ok))
	}

	return ok, nil
}