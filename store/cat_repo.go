package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type CatRepo struct {
	conn Conn
}

func (c *CatRepo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (DboSlice, error) {
	var key Track
	var app App
	var keys []DBO

	_, err := c.conn.QueryFunc(
		ctx,
		sql,
		params,
		[]interface{}{
			&key.Id, &key.BundleId, &key.Type, &key.Place, &key.Date,
			&app.Id, &app.Bundle, &app.Category, &app.DeveloperId, &app.Developer, &app.Geo,
			&app.StartAt, &app.Period,
		},
		func(row pgx.QueryFuncRow) error {
			key.App = app
			keys = append(keys, key)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, pgx.ErrNoRows
	}

	return keys, nil
}

func (c *CatRepo) ByBundleId(ctx context.Context, bundleId int) (DboSlice, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking CAT inner join app_tracking APP on CAT.bundleid = APP.id  where CAT.bundleid = $1",
		bundleId,
	)
}

func (c *CatRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (DboSlice, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking CAT inner join app_tracking APP on CAT.bundleid = APP.id where CAT.bundleid = $1 and CAT.date >= $2 and CAT.date <= $3",
		bundleId, start, end,
	)
}

func (c *CatRepo) LastUpdates(ctx context.Context, bundleId, count int) (DboSlice, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking CAT inner join app_tracking APP on CAT.bundleid = APP.id where CAT.bundleid = $1 order by CAT.id desc limit $2",
		bundleId, count,
	)
}

func NewCat(conn Conn) *CatRepo {
	return &CatRepo{
		conn: conn,
	}
}



