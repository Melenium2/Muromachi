package testhelpers

import (
	"Muromachi/config"
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strings"
	"time"
)

// Create instance of real database from local config
//
// Also return cleaner func for truncate data from tables
func RealDb() (*pgxpool.Pool, func(names ...string)) {
	cfg := config.New("../../config/dev.yml")
	url, err := connector.ConnectionUrl(cfg.Database)
	if err != nil {
		panic(err)
	}

	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		panic(err)
	}

	if err = connector.InitSchema(conn, "../../config/schema.sql"); err != nil {
		panic(err)
	}
	return conn, func(names ...string) {
		_, err = conn.Exec(context.Background(), fmt.Sprintf("truncate table %s CASCADE", strings.Join(names, ",")))
		if err != nil {
			log.Print(err)
		}
	}
}

// Insert new apprepo to app_tracking table in test database
func AddNewApp(conn *pgxpool.Pool, ctx context.Context, app entities.App) (int, error) {
	row := conn.QueryRow(
		ctx,
		fmt.Sprint("insert into app_tracking (bundle, category, developerId, developer, geo, startAt, period)  values ($1, $2, $3, $4, $5, $6, $7) returning id"),
		app.Bundle, app.Category, app.DeveloperId, app.Developer, app.Geo, app.StartAt, app.Period,
	)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// Insert new meta to meta_tracking table in test database
func AddNewMeta(conn *pgxpool.Pool, ctx context.Context, meta entities.Meta) (int, error) {
	values := "(bundleId, title, price, picture, screenshots," +
		" rating, reviewCount, ratingHistogram, description," +
		" shortDescription, recentChanges, releaseDate, lastUpdateDate, appSize," +
		" installs, version, androidVersion, contentRating, devContacts," +
		" privacyPolicy, date)"
	row := conn.QueryRow(
		ctx,
		fmt.Sprintf("insert into meta_tracking %s values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19::developerContacts, $20, $21) returning id", values),
		meta.BundleId,
		meta.Title,
		meta.Price,
		meta.Picture,
		meta.Screenshots,
		meta.Rating,
		meta.ReviewCount,
		meta.RatingHistogram,
		meta.Description,
		meta.ShortDescription,
		meta.RecentChanges,
		meta.ReleaseDate,
		meta.LastUpdateDate,
		meta.AppSize,
		meta.Installs,
		meta.Version,
		meta.AndroidVersion,
		meta.ContentRating,
		meta.DeveloperContacts,
		meta.PrivacyPolicy,
		meta.Date,
	)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// Add new track to keyword or category table in test database
func AddNewTrack(conn *pgxpool.Pool, ctx context.Context, track entities.Track, table string) (int, error) {
	row := conn.QueryRow(
		ctx,
		fmt.Sprintf("insert into %s (bundleId, type, place, date) values ($1, $2, $3, $4) returning id", table),
		track.BundleId,
		track.Type,
		track.Place,
		track.Date,
	)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// Return new meta struct for tests
func MetaStruct(bundleId int) entities.Meta {
	t1, _ := time.Parse("2006-01-02", "2021-01-18")
	return entities.Meta{
		BundleId:         bundleId,
		Title:            "Im title",
		Price:            "",
		Picture:          "http://picture",
		Screenshots:      []string{"http://picture", "http://picture1", "http://picture"},
		Rating:           "4.6+",
		ReviewCount:      "1002323",
		RatingHistogram:  []string{"1", "2", "3", "4", "5"},
		Description:      "some description of apprepo",
		ShortDescription: "some short description",
		RecentChanges:    "last changes",
		ReleaseDate:      "2020-01-01",
		LastUpdateDate:   "2020-03-03",
		AppSize:          "90MB+",
		Installs:         "1000000+",
		Version:          "v1.3.12",
		AndroidVersion:   "9.0",
		ContentRating:    "18+",
		DeveloperContacts: entities.DeveloperContacts{
			Email:    "email@email.com",
			Contacts: "virginia",
		},
		PrivacyPolicy: "http://privacypolicy.com",
		Date:          t1.AddDate(0, 0, 2),
	}
}

// Return new track struct for tests
func TrackStruct(bundleId int, t string) entities.Track {
	t1, _ := time.Parse("2006-01-02", "2021-01-18")
	return entities.Track{
		BundleId: bundleId,
		Type:     t,
		Date:     t1,
		Place:    19,
	}
}

