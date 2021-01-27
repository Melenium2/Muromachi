package store_test

import (
	"Muromachi/store"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAppRepo_ByBundleId_ShouldReturnSliceOfApps_Mock(t *testing.T) {
	conn := mockAppConnection{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 10)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var app store.App
	assert.NoError(t, dboSlice[0].To(&app))
	assert.Equal(t, "FINANCE", app.Category)
}

func TestAppRepo_ByBundleId_ShouldReturnSliceOfApps(t *testing.T) {
	conn, cleaner := RealDb()
	defer cleaner("app_tracking")
	repo := store.NewApp(conn)
	ctx := context.Background()

	app := store.App{
		Bundle:      "com.test.hello",
		Category:    "FINANCE",
		DeveloperId: "imdevid",
		Developer:   "invalid",
		Geo:         "ru_ru",
		StartAt:     time.Now(),
		Period:      31,
	}
	id, err := AddNewApp(conn, ctx, app)
	assert.NoError(t, err)

	dbo, err := repo.ByBundleId(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, dbo)

	var appFromDb store.App
	assert.NoError(t, dbo[0].To(&appFromDb))
	assert.Equal(t, id, appFromDb.Id)
}

func TestAppRepo_LastUpdates_ShouldReturnAllApplicationWithinGivenInterval_Mock(t *testing.T) {
	conn := mockAppConnection{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	timestamp, _ := time.Parse("2006-01-01", "2020-01-01")
	timestamp = timestamp.Add(time.Hour * 25)
	nextTimestamp := timestamp.AddDate(0, 1, 0)

	dboSlice, err := repo.TimeRange(ctx, 10, timestamp, nextTimestamp)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var app store.App
	var secondApp store.App
	assert.NoError(t, dboSlice[0].To(&app))
	assert.NoError(t, dboSlice[1].To(&secondApp))

	assert.True(t, timestamp.Equal(app.StartAt))
	assert.True(t, nextTimestamp.After(secondApp.StartAt))
}

func TestAppRepo_LastUpdates_ShouldReturnErrorBecauseTheFuncNotAllowedInThisTable(t *testing.T) {
	conn, _ := RealDb()
	repo := store.NewApp(conn)
	ctx := context.Background()

	_, err := repo.LastUpdates(ctx, 10, 10)
	assert.Error(t, err)
}

func TestAppRepo_LastUpdates_ShouldReturnErrorBecauseThisTableHasNotInfo_Mock(t *testing.T) {
	conn := mockAppConnection{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	_, err := repo.LastUpdates(ctx, 10, 1)
	assert.Error(t, err)
}

func TestAppRepo_ByBundleId_ShouldReturnErrNoRows_Mock(t *testing.T) {
	conn := mockAppConnectionErrors{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	_, err := repo.ByBundleId(ctx, 10)
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
}

func TestAppRepo_ByBundleIdShouldReturnErrNoRows_Mock(t *testing.T) {
	conn := mockAppConnectionErrors{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	_, err := repo.TimeRange(ctx, 10, time.Now(), time.Now().AddDate(0, 0, 1))
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
}
