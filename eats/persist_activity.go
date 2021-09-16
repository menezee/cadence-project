package eats

import (
	"context"
	"errors"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func PersistActivity(ctx context.Context, toBeSaved map[string]interface{}, workflowId string) (bool, error) {
	activity.GetLogger(ctx).Info("Persist Activity started.", zap.Any("toBeSaved", toBeSaved), zap.String("workflowId", workflowId))

	persistService := NewPersistTransactionService()
	ok := persistService.Persist(toBeSaved, workflowId)

	if ok {
		activity.GetLogger(ctx).Info("Persist Activity completed.", zap.Bool("success:", ok))
		return ok, nil
	}
	return ok, errors.New("payment error")
}
