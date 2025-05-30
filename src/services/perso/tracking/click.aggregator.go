package tracking

import (
	"github.com/Menschomat/bly.li/shared/mongo"
)

func AggregateClicks() {
	logger.Info("Calling click aggregator...")
	mongo.RunClickAggregation()
}
