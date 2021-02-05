package metastore_test

import (
	"Muromachi/config"
	"Muromachi/store/entities"
	"Muromachi/store/testhelpers"
	"Muromachi/store/tracking/metastore"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMetaRepo_ByBundleId_ShouldReturnSomeApps_Mock(t *testing.T) {
	conn := mockMetaConnection{}
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 12)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 3, len(dboSlice))

	var app entities.Meta
	assert.NoError(t, dboSlice[0].To(&app))
	assert.Equal(t, "Im title", app.Title)
	assert.Equal(t, 3, len(app.Screenshots))
	assert.Equal(t, "123", app.App.Bundle)
	assert.NotEmpty(t, app.DeveloperContacts.Contacts)
	assert.NotEmpty(t, app.DeveloperContacts.Email)
}

func TestMetaRepo_ByBundleId_ShouldReturnSomeApps(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("app_tracking, meta_tracking")
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	bundleId, err := testhelpers.AddNewApp(conn, ctx, entities.App{
		Bundle: "123",
	})

	meta := testhelpers.MetaStruct(bundleId)
	for i := 0; i < 3; i++ {
		_, err := testhelpers.AddNewMeta(conn, ctx, meta)
		assert.NoError(t, err)
		meta.Date = meta.Date.AddDate(0, 0, 1 + i)
	}

	dboSlice, err := repo.ByBundleId(ctx, bundleId)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var app entities.Meta
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&app))
		assert.Equal(t, bundleId, app.BundleId)
		assert.Equal(t, "123", app.App.Bundle)
	}
}

func TestMetaRepo_TimeRange_ShouldReturnAppsWithGivenTimeRange_Mock(t *testing.T) {
	conn := mockMetaConnection{}
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	t1, _ := time.Parse("2006-01-02", "2021-01-18")
	t2 := t1.AddDate(0, 0, 7)
	dboSlice, err := repo.TimeRange(ctx, 12, t1, t2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var app, lastApp entities.Meta
	assert.NoError(t, dboSlice[0].To(&app))
	assert.NoError(t, dboSlice[2].To(&lastApp))

	assert.True(t, app.Date.After(t1))
	assert.True(t, lastApp.Date.Before(t2))
	assert.Equal(t, "123", lastApp.App.Bundle)
}

func TestMetaRepo_TimeRange_ShouldReturnAppsWithGivenTimeRange(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("app_tracking, meta_tracking")
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	bundleId, err := testhelpers.AddNewApp(conn, ctx, entities.App{
		Bundle: "123",
	})

	meta := testhelpers.MetaStruct(bundleId)
	for i := 0; i < 3; i++ {
		_, err := testhelpers.AddNewMeta(conn, ctx, meta)
		assert.NoError(t, err)
		meta.Date = meta.Date.AddDate(0, 0, 1 + i)
	}
	t1 := meta.Date.AddDate(0, 0, -10)
	t2 := meta.Date.AddDate(0, 0, 4)
	dboSlice, err := repo.TimeRange(ctx, bundleId, t1, t2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var app entities.Meta
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&app))
		assert.Equal(t, bundleId, app.BundleId)
		assert.Equal(t, "123", app.App.Bundle)
		assert.True(t, app.Date.After(t1) && (app.Date.Before(t2) || app.Date.Equal(t2)) )
	}
}

func TestMetaRepo_LastUpdates_ShouldReturnLastNApps_Mock(t *testing.T) {
	conn := mockMetaConnection{}
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.LastUpdates(ctx, 12, 3)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 3, len(dboSlice))
}

func TestMetaRepo_LastUpdates_ShouldReturnLastNApps(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("app_tracking, meta_tracking")
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	bundleId, err := testhelpers.AddNewApp(conn, ctx, entities.App{
		Bundle: "123",
	})

	meta := testhelpers.MetaStruct(bundleId)
	for i := 0; i < 4; i++ {
		_, err := testhelpers.AddNewMeta(conn, ctx, meta)
		assert.NoError(t, err)
		meta.Date = meta.Date.AddDate(0, 0, 1 + i)
	}

	dboSlice, err := repo.LastUpdates(ctx, bundleId, 2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 2, len(dboSlice))

	var app entities.Meta
	id := 1000
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&app))
		assert.Equal(t, bundleId, app.BundleId)
		assert.Equal(t, "123", app.App.Bundle)
		assert.Greater(t, id, app.Id)
		id = app.Id
	}
}

func TestMetaRepo_ByBundleId_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockMetaConnectionErrors{}
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 12)
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
	assert.Nil(t, dboSlice)
}

func TestMetaRepo_ByBundleId_ShouldReturnErrorIfNoRows(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, _ := testhelpers.RealDb(cfg.Database)
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 12)
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
	assert.Nil(t, dboSlice)
}

func TestMetaRepo_TimeRange_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockMetaConnectionErrors{}
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.TimeRange(ctx, 12, time.Now(), time.Now())
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
	assert.Nil(t, dboSlice)
}

func TestMetaRepo_TimeRange_ShouldReturnErrorIfNoRows(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, _ := testhelpers.RealDb(cfg.Database)
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.TimeRange(ctx, 12, time.Now(), time.Now())
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
	assert.Nil(t, dboSlice)
}

func TestMetaRepo_LastUpdate_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockMetaConnectionErrors{}
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.LastUpdates(ctx, 12, 1)
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
	assert.Nil(t, dboSlice)
}

func TestMetaRepo_LastUpdate_ShouldReturnErrorIfNoRows(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, _ := testhelpers.RealDb(cfg.Database)
	repo := metastore.Repo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.LastUpdates(ctx, 12, 1)
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
	assert.Nil(t, dboSlice)
}

