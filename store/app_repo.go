package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type AppRepo struct {
	conn Conn
}

func (a *AppRepo) ById(ctx context.Context, id int) (DBO, error) {
	var app App
	err := a.conn.QueryRow(
		ctx,
		"select * from app_tracking where id = $1",
		id,
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

	return app, nil
}

func (a *AppRepo) ByBundleId(ctx context.Context, bundleId int) ([]DBO, error) {
	var app App
	var apps []DBO
	_, err := a.conn.QueryFunc(
		ctx,
		"select * from app_tracking where bundleid = $1",
		[]interface{}{bundleId},
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

func (a *AppRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) ([]DBO, error) {
	panic("implement me")
}

func (a *AppRepo) LastUpdates(ctx context.Context, bundleId, count int) ([]DBO, error) {
	panic("implement me")
}

func NewApp(conn Conn) *AppRepo {
	return &AppRepo{
		conn: conn,
	}
}
