package eats

import (
	"context"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func ETAActivity(ctx context.Context) (map[string]interface{}, error) {
	activity.GetLogger(ctx).Info("ETA Activity started.")

	//Generate ETA:
	rand.Seed(time.Now().UnixNano())
	min := 2
	max := 25
	etaRandom := rand.Intn(max - min + 1) + min
	activity.GetLogger(ctx).Info("ETA Random.")

	//CourierName
	possibleCouriers := []string{"Piera", "Erich", "Demontie", "Gisela", "Keiko"}
	courierName := possibleCouriers[rand.Intn(4)]

	response := map[string]interface{}{"ETA": etaRandom, "CourierName": courierName}

	activity.GetLogger(ctx).Info("ETA Activity completed.", zap.Any("response:", response))

	return response, nil
}
