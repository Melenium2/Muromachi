package store_test

import (
	"Muromachi/config"
	"Muromachi/store"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"log"
	"strings"
	"time"
)

// Create instance of real database from local config
//
// Also return cleaner func for truncate data from tables
func RealDb() (*pgx.Conn, func(names ...string)) {
	c := config.New("../config/dev.yml")
	url, err := store.ConnectionUrl(c.Database)
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		panic(err)
	}

	if err := store.InitSchema(conn, "../config/schema.sql"); err != nil {
		panic(err)
	}

	return conn, func(names ...string) {
		_, err := conn.Exec(context.Background(), fmt.Sprintf("truncate table %s CASCADE", strings.Join(names, ",")))
		if err != nil {
			log.Print(err)
		}
	}
}

// Insert new app to app_tracking table in test database
func AddNewApp(conn *pgx.Conn, ctx context.Context, app store.App) (int, error) {
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
func AddNewMeta(conn *pgx.Conn, ctx context.Context, meta store.Meta) (int, error) {
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

// Return new meta struct for tests
func MetaStruct(bundleId int) store.Meta {
	t1, _ := time.Parse("2006-01-02", "2021-01-18")
	return store.Meta{
		BundleId:         bundleId,
		Title:            "Im title",
		Price:            "",
		Picture:          "http://picture",
		Screenshots:      []string{"http://picture", "http://picture1", "http://picture"},
		Rating:           "4.6+",
		ReviewCount:      "1002323",
		RatingHistogram:  []string{"1", "2", "3", "4", "5"},
		Description:      "some description of app",
		ShortDescription: "some short description",
		RecentChanges:    "last changes",
		ReleaseDate:      "2020-01-01",
		LastUpdateDate:   "2020-03-03",
		AppSize:          "90MB+",
		Installs:         "1000000+",
		Version:          "v1.3.12",
		AndroidVersion:   "9.0",
		ContentRating:    "18+",
		DeveloperContacts: store.DeveloperContacts{
			Email:    "email@email.com",
			Contacts: "virginia",
		},
		PrivacyPolicy: "http://privacypolicy.com",
		Date:          t1.AddDate(0, 0, 2),
	}
}

// Mock connection with errors (App table)
type mockAppConnectionErrors struct {
}

func (m mockAppConnectionErrors) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRowError{}
}

func (m mockAppConnectionErrors) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, pgx.ErrNoRows
}

type mockRowError struct {
}

func (mre mockRowError) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}

func (mre mockRowError) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mre mockRowError) RawValues() [][]byte {
	return nil
}

// Mock connection with successful returned objects (App table)
type mockAppConnection struct {
}

func (m mockAppConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	t, _ := time.Parse("2006-01-01", "2020-01-01")
	return mockRow{
		Id:          1,
		Bundle:      "com.test",
		Category:    "FINANCE",
		DeveloperId: "92834848476158744",
		Developer:   "Random",
		Geo:         "ru_ru",
		StartAt:     t,
		Period:      30,
	}
}

func (m mockAppConnection) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	t, _ := time.Parse("2006-01-02", "2020-01-01")
	for i := 0; i < 2; i++ {
		t = t.Add(time.Hour * 25)
		err := mockRow{
			Id:          i + 1,
			Bundle:      "com.test",
			Category:    "FINANCE",
			DeveloperId: "92834848476158744",
			Developer:   "Random",
			Geo:         "ru_ru",
			StartAt:     t,
			Period:      30,
		}.Scan(scans...)
		if err != nil {
			return nil, err
		}
		err = f(mockRow{})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type mockRow struct {
	Id          int
	Bundle      string
	Category    string
	DeveloperId string
	Developer   string
	Geo         string
	StartAt     time.Time
	Period      uint32
}

func (mr mockRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockRow) RawValues() [][]byte {
	return nil
}

func (mr mockRow) Scan(dest ...interface{}) error {
	Id := dest[0].(*int)
	Bundle := dest[1].(*string)
	Category := dest[2].(*string)
	DeveloperId := dest[3].(*string)
	Developer := dest[4].(*string)
	Geo := dest[5].(*string)
	StartAt := dest[6].(*time.Time)
	Period := dest[7].(*uint32)

	*Id = mr.Id
	*Bundle = mr.Bundle
	*Category = mr.Category
	*DeveloperId = mr.DeveloperId
	*Developer = mr.Developer
	*Geo = mr.Geo
	*StartAt = mr.StartAt
	*Period = mr.Period

	return nil
}

// Mock connection with errors (Meta table)
type mockMetaConnectionErrors struct {
}

func (m mockMetaConnectionErrors) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockMetaErrorRow{}
}

func (m mockMetaConnectionErrors) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, pgx.ErrNoRows
}

type mockMetaErrorRow struct {
}

func (mr mockMetaErrorRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockMetaErrorRow) RawValues() [][]byte {
	return nil
}

func (mr mockMetaErrorRow) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}

// Mock connection with successful returned objects (Meta table)
type mockMetaConnection struct {
}

func (m mockMetaConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	t, _ := time.Parse("2006-01-02", "2021-01-19")
	return mockMetaRow{
		Id:               10,
		BundleId:         12,
		Title:            "Im title",
		Price:            "",
		Picture:          "http://picture",
		Screenshots:      []string{"http://picture", "http://picture1", "http://picture"},
		Rating:           "4.6+",
		ReviewCount:      "1002323",
		RatingHistogram:  []string{"1", "2", "3", "4", "5"},
		Description:      "some description of app",
		ShortDescription: "some short description",
		RecentChanges:    "last changes",
		ReleaseDate:      "2020-01-01",
		LastUpdateDate:   "2020-03-03",
		AppSize:          "90MB+",
		Installs:         "1000000+",
		Version:          "v1.3.12",
		AndroidVersion:   "9.0",
		ContentRating:    "18+",
		DeveloperContacts: store.DeveloperContacts{
			Email:    "email@email.com",
			Contacts: "virginia",
		},
		PrivacyPolicy: "http://privacypolicy.com",
		Date:          t,
	}
}

func (m mockMetaConnection) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	t, _ := time.Parse("2006-01-02", "2021-01-19")
	for i := 0; i < 3; i++ {
		_ = mockMetaRow{
			Id:               10 + i,
			BundleId:         12,
			Title:            "Im title",
			Price:            "",
			Picture:          "http://picture",
			Screenshots:      []string{"http://picture", "http://picture1", "http://picture"},
			Rating:           "4.6+",
			ReviewCount:      "1002323",
			RatingHistogram:  []string{"1", "2", "3", "4", "5"},
			Description:      "some description of app",
			ShortDescription: "some short description",
			RecentChanges:    "last changes",
			ReleaseDate:      "2020-01-01",
			LastUpdateDate:   "2020-03-03",
			AppSize:          "90MB+",
			Installs:         "1000000+",
			Version:          "v1.3.12",
			AndroidVersion:   "9.0",
			ContentRating:    "18+",
			DeveloperContacts: store.DeveloperContacts{
				Email:    "email@email.com",
				Contacts: "virginia",
			},
			PrivacyPolicy: "http://privacypolicy.com",
			Date:          t,
		}.Scan(scans...)

		_ = f(mockMetaRow{})

		t = t.AddDate(0, 0, 1)
	}

	return nil, nil
}

type mockMetaRow struct {
	Id                int                     `json:"-"`
	BundleId          int                     `json:"-"`
	Title             string                  `json:"title" db:"title"`
	Price             string                  `json:"price" db:"price"`
	Picture           string                  `json:"picture" db:"picture"`
	Screenshots       []string                `json:"screenshots" db:"screenshots"`
	Rating            string                  `json:"rating" db:"rating"`
	ReviewCount       string                  `json:"reviewCount" db:"review_count"`
	RatingHistogram   []string                `json:"ratingHistogram" db:"rating_histogram"`
	Description       string                  `json:"description" db:"description"`
	ShortDescription  string                  `json:"shortDescription" db:"short_description"`
	RecentChanges     string                  `json:"recentChanges" db:"recent_changes"`
	ReleaseDate       string                  `json:"releaseDate" db:"release_date"`
	LastUpdateDate    string                  `json:"lastUpdateDate" db:"last_update_date"`
	AppSize           string                  `json:"appSize" db:"app_size"`
	Installs          string                  `json:"installs" db:"installs"`
	Version           string                  `json:"version" db:"version"`
	AndroidVersion    string                  `json:"androidVersion" db:"android_version"`
	ContentRating     string                  `json:"contentRating" db:"content_rating"`
	DeveloperContacts store.DeveloperContacts `json:"developerContacts" db:"developer_contacts"`
	PrivacyPolicy     string                  `json:"privacyPolicy,omitempty"`
	Date              time.Time               `json:"date,omitempty"`
}

func (mr mockMetaRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockMetaRow) RawValues() [][]byte {
	return nil
}

func (mr mockMetaRow) Scan(dest ...interface{}) error {
	Id := dest[0].(*int)
	BundleId := dest[1].(*int)
	Title := dest[2].(*string)
	Price := dest[3].(*string)
	Picture := dest[4].(*string)
	Screenshots := dest[5].(*[]string)
	Rating := dest[6].(*string)
	ReviewCount := dest[7].(*string)
	RatingHistogram := dest[8].(*[]string)
	Description := dest[9].(*string)
	ShortDescription := dest[10].(*string)
	RecentChanges := dest[11].(*string)
	ReleaseDate := dest[12].(*string)
	LastUpdateDate := dest[13].(*string)
	AppSize := dest[14].(*string)
	Installs := dest[15].(*string)
	Version := dest[16].(*string)
	AndroidVersion := dest[17].(*string)
	ContentRating := dest[18].(*string)
	DeveloperContacts := dest[19].(*store.DeveloperContacts)
	PrivacyPolicy := dest[20].(*string)
	Date := dest[21].(*time.Time)

	*Id = mr.Id
	*BundleId = mr.BundleId
	*Title = mr.Title
	*Price = mr.Price
	*Picture = mr.Picture
	*Screenshots = mr.Screenshots
	*Rating = mr.Rating
	*ReviewCount = mr.ReviewCount
	*RatingHistogram = mr.RatingHistogram
	*Description = mr.Description
	*ShortDescription = mr.ShortDescription
	*RecentChanges = mr.RecentChanges
	*ReleaseDate = mr.ReleaseDate
	*LastUpdateDate = mr.LastUpdateDate
	*AppSize = mr.AppSize
	*Installs = mr.Installs
	*Version = mr.Version
	*AndroidVersion = mr.AndroidVersion
	*ContentRating = mr.ContentRating
	*DeveloperContacts = mr.DeveloperContacts
	*PrivacyPolicy = mr.PrivacyPolicy
	*Date = mr.Date

	return nil
}
