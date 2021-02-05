package trackstore

import (
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

// Struct for holds connection with db
type CatRepo struct {
	Conn connector.Conn
}

// Making database queries
func (c *CatRepo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (entities.DboSlice, error) {
	var key entities.Track
	var app entities.App
	var keys []entities.DBO

	_, err := c.Conn.QueryFunc(
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

// Return DboSlice with entities.Track with bundleId which is equal to given id
func (c *CatRepo) ByBundleId(ctx context.Context, bundleId int) (entities.DboSlice, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking CAT inner join app_tracking APP on CAT.bundleid = APP.id  where CAT.bundleid = $1",
		bundleId,
	)
}

// Return DboSlice with given bundle id and within time range from start to end
func (c *CatRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (entities.DboSlice, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking CAT inner join app_tracking APP on CAT.bundleid = APP.id where CAT.bundleid = $1 and CAT.date >= $2 and CAT.date <= $3",
		bundleId, start, end,
	)
}

// Get last n updates of categories with bundle id equals given bundle id
func (c *CatRepo) LastUpdates(ctx context.Context, bundleId, count int) (entities.DboSlice, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking CAT inner join app_tracking APP on CAT.bundleid = APP.id where CAT.bundleid = $1 order by CAT.id desc limit $2",
		bundleId, count,
	)
}



