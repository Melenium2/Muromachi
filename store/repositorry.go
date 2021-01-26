package store

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

type TrackingRepo interface {
	ProducerFunc(ctx context.Context, sql string, params ...interface{}) ([]DBO, error)
	ByBundleId(ctx context.Context, bundleId int) ([]DBO, error)
	TimeRange(ctx context.Context, bundleId int, start, end time.Time) ([]DBO, error)
	LastUpdates(ctx context.Context, bundleId, count int) ([]DBO, error)
}

type Conn interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
}