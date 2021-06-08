package eats

import (
	"context"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func CreditActivity(ctx context.Context, amount int) (bool, error) {
	activity.GetLogger(ctx).Info("CreditActivity started.", zap.Int("Amount", amount))
	accountTransactionService := NewAccountTransactionService()
	ok := accountTransactionService.Credit(amount)

	if ok {
		activity.GetLogger(ctx).Info("CreditActivity completed.", zap.Bool("Succeeded", ok))
	}

	return ok, nil
}