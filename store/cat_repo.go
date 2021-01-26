package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type CatRepo struct {
	conn Conn
}

func (c *CatRepo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) ([]DBO, error) {
	var key Track
	var keys []DBO

	_, err := c.conn.QueryFunc(
		ctx,
		sql,
		params,
		[]interface{}{
			&key.Id, &key.BundleId, &key.Type, &key.Place, &key.Date,
		},
		func(row pgx.QueryFuncRow) error {
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

func (c *CatRepo) ByBundleId(ctx context.Context, bundleId int) ([]DBO, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking where bundleid = $1",
		bundleId,
	)
}

func (c *CatRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) ([]DBO, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking where bundleid = $1 and date >= $2 and date <= $3",
		bundleId, start, end,
	)
}

func (c *CatRepo) LastUpdates(ctx context.Context, bundleId, count int) ([]DBO, error) {
	return c.ProducerFunc(
		ctx,
		"select * from category_tracking where bundleid = $1 order by id desc limit $2",
		bundleId, count,
	)
}

func NewCat(conn Conn) *CatRepo {
	return &CatRepo{
		conn: conn,
	}
}



