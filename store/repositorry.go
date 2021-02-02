package store

import (
	"Muromachi/store/entities"
	"context"
	"time"
)



type TrackingRepo interface {
	ProducerFunc(ctx context.Context, sql string, params ...interface{}) (entities.DboSlice, error)
	ByBundleId(ctx context.Context, bundleId int) (entities.DboSlice, error)
	TimeRange(ctx context.Context, bundleId int, start, end time.Time) (entities.DboSlice, error)
	LastUpdates(ctx context.Context, bundleId, count int) (entities.DboSlice, error)
}
