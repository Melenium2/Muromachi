package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type MetaRepo struct {
	conn Conn
}

func (m *MetaRepo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (DboSlice, error) {
	var app Meta
	var apps []DBO

	_, err := m.conn.QueryFunc(
		ctx,
		sql,
		append([]interface{}{pgx.QueryResultFormats{pgx.BinaryFormatCode}}, params...),
		[]interface{}{
			&app.Id, &app.BundleId, &app.Title, &app.Price, &app.Picture,
			&app.Screenshots, &app.Rating, &app.ReviewCount, &app.RatingHistogram,
			&app.Description, &app.ShortDescription, &app.RecentChanges, &app.ReleaseDate,
			&app.LastUpdateDate, &app.AppSize, &app.Installs, &app.Version, &app.AndroidVersion,
			&app.ContentRating, &app.DeveloperContacts, &app.PrivacyPolicy, &app.Date,
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

func (m *MetaRepo) ByBundleId(ctx context.Context, bundleId int) (DboSlice, error) {
	return m.ProducerFunc(
		ctx,
		"select * from meta_tracking where bundleid = $1",
		bundleId,
	)
}

func (m *MetaRepo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (DboSlice, error) {
	return m.ProducerFunc(
		ctx,
		"select * from meta_tracking where bundleid = $1 and date >= $2 and date <= $3",
		bundleId, start, end,
	)
}

func (m *MetaRepo) LastUpdates(ctx context.Context, bundleId, count int) (DboSlice, error) {
	return m.ProducerFunc(
		ctx,
		"select * from meta_tracking where bundleid = $1 order by id desc limit $2",
		bundleId, count,
	)
}

func NewMeta(conn Conn) *MetaRepo {
	return &MetaRepo{
		conn: conn,
	}
}
