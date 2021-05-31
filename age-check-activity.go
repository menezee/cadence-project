package main

import (
	"context"
	"fmt"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func AgeCheckActivity(ctx context.Context, age int) (string, error) {
	ageAsStr := fmt.Sprintf("%d", age)
	activity.GetLogger(ctx).Info("AgeCheckActivity called.", zap.String("Age", ageAsStr))

	if age >= 18 {
		return ageAsStr + " is an adult!", nil
	}

	return ageAsStr + " is underage", nil
}
