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

func (a *AppRepo) ByBundleId(ctx context.Context, bundleId int) ([]DBO, error) {
	var app App
	err := a.conn.QueryRow(
		ctx,
		"select * from app_tracking where id = $1",
		bundleId,
	).
		Scan(
			&app.Id,
			&app.Bundle,
			&app.Category,
			&app.DeveloperId,
			&app.Developer,
			&app.Geo,
			&app.StartAt,
			&app.Period,
		)

	if err != nil {
		return nil, err
	}

	return []DBO{app}, nil
}

func (a *AppRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) ([]DBO, error) {
	var app App
	var apps []DBO
	_, err := a.conn.QueryFunc(
		ctx,
		"select * from app_tracking where bundleid = $1 and startat >= $2 and startat <= $3",
		[]interface{}{bundleId, start, end},
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

	return apps, nil
}

func (a *AppRepo) LastUpdates(_ context.Context, _, _ int) ([]DBO, error) {
	return nil, fmt.Errorf("%s", "no last updates in this table")
}

func NewApp(conn Conn) *AppRepo {
	return &AppRepo{
		conn: conn,
	}
}
