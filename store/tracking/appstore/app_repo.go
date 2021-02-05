package appstore

import (
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

// Struct for holds connection with db
type Repo struct {
	Conn connector.Conn
}

// Making database queries
func (a *Repo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (entities.DboSlice, error) {
	var app entities.App
	var apps []entities.DBO

	_, err := a.Conn.QueryFunc(
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

// Return DboSlice with entities.App with bundleId which is equal to given id
func (a *Repo) ByBundleId(ctx context.Context, bundleId int) (entities.DboSlice, error) {
	return a.ProducerFunc(
		ctx,
		"select * from app_tracking where id = $1",
		bundleId,
	)
}

// Return DboSlice with given bundle id and within time range from start to end
func (a *Repo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (entities.DboSlice, error) {
	return a.ProducerFunc(
		ctx,
		"select * from app_tracking where bundleid = $1 and startat >= $2 and startat <= $3",
		bundleId, start, end,
	)
}

// Do nothing here
func (a *Repo) LastUpdates(_ context.Context, _, _ int) (entities.DboSlice, error) {
	return nil, fmt.Errorf("%s", "no last updates in this table")
}
