package store

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

type Conn interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
}

type TrackingRepo interface {
	ProducerFunc(ctx context.Context, sql string, params ...interface{}) (DboSlice, error)
	ByBundleId(ctx context.Context, bundleId int) (DboSlice, error)
	TimeRange(ctx context.Context, bundleId int, start, end time.Time) (DboSlice, error)
	LastUpdates(ctx context.Context, bundleId, count int) (DboSlice, error)
}

type UsersRepo interface {
	Create(ctx context.Context, user User) (User, error)
	Approve(ctx context.Context, clientId, clientSecret string) (User, error)
}

type BlackList interface {
	AddBlock()
	CheckBlock()
}

type RefreshSessions interface {
	New(ctx context.Context, session Session) (Session, error)
	Get(ctx context.Context, token string) (Session, error)
	Remove(ctx context.Context, token string) (Session, error)
	RemoveBatch(ctx context.Context, sessionid ...int) error
	UserSessions(ctx context.Context, userId int) ([]Session, error)
}

type Sessions interface {
	BlackList
	RefreshSessions
}
