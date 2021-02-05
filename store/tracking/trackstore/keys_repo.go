package trackstore

import (
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

// Struct for holds connection with db
type KeysRepo struct {
	Conn connector.Conn
}

// Making database queries
func (k *KeysRepo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (entities.DboSlice, error) {
	var (
		key  entities.Track
		app  entities.App
		keys []entities.DBO
	)

	_, err := k.Conn.QueryFunc(
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
func (k *KeysRepo) ByBundleId(ctx context.Context, bundleId int) (entities.DboSlice, error) {
	return k.ProducerFunc(
		ctx,
		"select * from keyword_tracking KEY inner join app_tracking APP on KEY.bundleid = APP.id where KEY.bundleid = $1",
		bundleId,
	)
}

// Return DboSlice with given bundle id and within time range from start to end
func (k *KeysRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (entities.DboSlice, error) {
	return k.ProducerFunc(
		ctx,
		"select * from keyword_tracking KEY inner join app_tracking APP on KEY.bundleid = APP.id where KEY.bundleid = $1 and KEY.date >= $2 and KEY.date <= $3",
		bundleId, start, end,
	)
}

// Get last n updates of app with bundle id equals given bundle id
func (k *KeysRepo) LastUpdates(ctx context.Context, bundleId, count int) (entities.DboSlice, error) {
	return k.ProducerFunc(
		ctx,
		"select * from keyword_tracking KEY inner join app_tracking APP on KEY.bundleid = APP.id where KEY.bundleid = $1 order by KEY.id desc limit $2",
		bundleId, count,
	)
}
