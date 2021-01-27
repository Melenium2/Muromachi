package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type KeysRepo struct {
	conn Conn
}

func (k *KeysRepo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (DboSlice, error) {
	var (
		key  Track
		app  App
		keys []DBO
	)

	_, err := k.conn.QueryFunc(
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

func (k *KeysRepo) ByBundleId(ctx context.Context, bundleId int) (DboSlice, error) {
	return k.ProducerFunc(
		ctx,
		"select * from keyword_tracking KEY inner join app_tracking APP on KEY.bundleid = APP.id where KEY.bundleid = $1",
		bundleId,
	)
}

func (k *KeysRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (DboSlice, error) {
	return k.ProducerFunc(
		ctx,
		"select * from keyword_tracking KEY inner join app_tracking APP on KEY.bundleid = APP.id where KEY.bundleid = $1 and KEY.date >= $2 and KEY.date <= $3",
		bundleId, start, end,
	)
}

func (k *KeysRepo) LastUpdates(ctx context.Context, bundleId, count int) (DboSlice, error) {
	return k.ProducerFunc(
		ctx,
		"select * from keyword_tracking KEY inner join app_tracking APP on KEY.bundleid = APP.id where KEY.bundleid = $1 order by KEY.id desc limit $2",
		bundleId, count,
	)
}

func NewKeys(conn Conn) *KeysRepo {
	return &KeysRepo{
		conn: conn,
	}
}
