package store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type AppRepo struct {
	conn Conn
}

func (a *AppRepo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (DboSlice, error) {
	var app App
	var apps []DBO

	_, err := a.conn.QueryFunc(
		ctx,
		sql,
		params,
		[]interface{}{
			&app.Id, &app.Bundle, &app.Category, &app.DeveloperId,
			&app.Developer, &app.Geo, &app.StartAt, &app.Period,
		},
		func(row pgx.QueryFuncRow) error {
			apps = append(apps, app)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		return nil, pgx.ErrNoRows
	}

	return apps, nil
}

func (a *AppRepo) ByBundleId(ctx context.Context, bundleId int) (DboSlice, error) {
	return a.ProducerFunc(
		ctx,
		"select * from app_tracking where id = $1",
		bundleId,
	)
}

func (a *AppRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (DboSlice, error) {
	return a.ProducerFunc(
		ctx,
		"select * from app_tracking where bundleid = $1 and startat >= $2 and startat <= $3",
		bundleId, start, end,
	)
}

func (a *AppRepo) LastUpdates(_ context.Context, _, _ int) (DboSlice, error) {
	return nil, fmt.Errorf("%s", "no last updates in this table")
}

func NewApp(conn Conn) *AppRepo {
	return &AppRepo{
		conn: conn,
	}
}
