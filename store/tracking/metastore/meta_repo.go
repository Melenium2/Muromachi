package metastore

import (
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

// Struct for holds connection with db
type Repo struct {
	Conn connector.Conn
}

// Making database queries
func (m *Repo) ProducerFunc(ctx context.Context, sql string, params ...interface{}) (entities.DboSlice, error) {
	var (
		meta entities.Meta
		app  entities.App
		apps []entities.DBO
	)

	_, err := m.Conn.QueryFunc(
		ctx,
		sql,
		append([]interface{}{pgx.QueryResultFormats{pgx.BinaryFormatCode}}, params...),
		[]interface{}{
			&meta.Id, &meta.BundleId, &meta.Title, &meta.Price, &meta.Picture,
			&meta.Screenshots, &meta.Rating, &meta.ReviewCount, &meta.RatingHistogram,
			&meta.Description, &meta.ShortDescription, &meta.RecentChanges, &meta.ReleaseDate,
			&meta.LastUpdateDate, &meta.AppSize, &meta.Installs, &meta.Version, &meta.AndroidVersion,
			&meta.ContentRating, &meta.DeveloperContacts, &meta.PrivacyPolicy, &meta.Date,
			&app.Id, &app.Bundle, &app.Category, &app.DeveloperId, &app.Developer, &app.Geo,
			&app.StartAt, &app.Period,
		},
		func(row pgx.QueryFuncRow) error {
			meta.App = app
			apps = append(apps, meta)
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

// Return DboSlice with entities.Meta with bundleId which is equal to given id
func (m *Repo) ByBundleId(ctx context.Context, bundleId int) (entities.DboSlice, error) {
	return m.ProducerFunc(
		ctx,
		"select * from meta_tracking inner join app_tracking APP on bundleid = APP.id  where bundleid = $1",
		bundleId,
	)
}

// Return DboSlice with given bundle id and within time range from start to end
func (m *Repo) TimeRange(ctx context.Context, bundleId int, start, end time.Time) (entities.DboSlice, error) {
	return m.ProducerFunc(
		ctx,
		"select * from meta_tracking inner join app_tracking APP on bundleid = APP.id  where bundleid = $1 and date >= $2 and date <= $3",
		bundleId, start, end,
	)
}

// Get last n updates of app with bundle id equals given bundle id
func (m *Repo) LastUpdates(ctx context.Context, bundleId, count int) (entities.DboSlice, error) {
	return m.ProducerFunc(
		ctx,
		"select * from meta_tracking META inner join app_tracking APP on bundleid = APP.id  where bundleid = $1 order by META.id desc limit $2",
		bundleId, count,
	)
}
