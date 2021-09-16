package eats

import (
	"context"
	"errors"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
	"log"
)

func PaymentActivity(ctx context.Context, amount int) (map[string]interface{}, error) {
	activity.GetLogger(ctx).Info("Payment Activity started.", zap.Int("Amount", amount))
	accountTransactionService := NewAccountTransactionService()
	ok := accountTransactionService.Debit(amount)

	if ok {
		response := map[string]interface{}{"PaymentApproved": ok}
		activity.GetLogger(ctx).Info("Payment Activity completed.", zap.Any("response:", response))
		log.Println("Payment Activity completed: ", response)
		return response, nil
	}
	return nil, errors.New("payment error")
}
