package main

import (
	"context"

	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

// SimpleActivity is a sample Cadence activity function that takes one parameter and
// returns a string containing the parameter value.
func SimpleActivity(ctx context.Context, value string) (string, error) {
	activity.GetLogger(ctx).Info("SimpleActivity called.", zap.String("Value", value))
	return "Processed: " + value, nil
}
