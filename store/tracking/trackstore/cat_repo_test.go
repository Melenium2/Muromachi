package trackstore_test

import (
	"Muromachi/config"
	"Muromachi/store/entities"
	"Muromachi/store/testhelpers"
	"Muromachi/store/tracking/trackstore"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCatRepo_ByBundleId_ShouldReturnApp_Mock(t *testing.T) {
	conn := mockTrackConnection{}
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 123)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var key entities.Track
	assert.NoError(t, dboSlice[0].To(&key))
	assert.Equal(t, "type", key.Type)
	assert.Equal(t, "123", key.App.Bundle)
}

func TestCatRepo_ByBundleId_ShouldReturnApp(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("app_tracking", "category_tracking")
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	bundleId, _ := testhelpers.AddNewApp(conn, ctx, entities.App{Bundle: "123"})
	track := testhelpers.TrackStruct(bundleId, "key")
	for i := 0; i < 4; i++ {
		_, _ = testhelpers.AddNewTrack(conn, ctx, track, "category_tracking")
	}

	dboSlice, err := repo.ByBundleId(ctx, bundleId)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 4, len(dboSlice))

	var key entities.Track
	assert.NoError(t, dboSlice[0].To(&key))
	assert.Equal(t, "key", key.Type)
	assert.Equal(t, "123", key.App.Bundle)
}

func TestCatRepo_TimeRange_ShouldReturnAppsInTimeRange_Mock(t *testing.T) {
	conn := mockTrackConnection{}
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	t1, _ := time.Parse("2006-01-02", "2021-01-18")
	t2 := t1.AddDate(0, 0, 10)
	dboSlice, err := repo.TimeRange(ctx, 123, t1, t2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var key entities.Track
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&key))
		assert.Equal(t, "type", key.Type)
		assert.Equal(t, "123", key.App.Bundle)
		assert.True(t, key.Date.After(t1) && key.Date.Before(t2))
	}
}

func TestCatRepo_TimeRange_ShouldReturnAppsInTimeRange(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("app_tracking", "category_tracking")
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	bundleId, _ := testhelpers.AddNewApp(conn, ctx, entities.App{Bundle: "123"})
	track := testhelpers.TrackStruct(bundleId, "key")
	t1 := track.Date.AddDate(0, 0, -1)
	for i := 0; i < 4; i++ {
		_, _ = testhelpers.AddNewTrack(conn, ctx, track, "category_tracking")
		track.Date = track.Date.AddDate(0, 0, 1)
	}
	t2 := track.Date.AddDate(0, 0, 1)

	dboSlice, err := repo.TimeRange(ctx, bundleId, t1, t2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 4, len(dboSlice))

	var key entities.Track
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&key))
		assert.Equal(t, "key", key.Type)
		assert.Equal(t, "123", key.App.Bundle)
		assert.True(t, key.Date.After(t1) && key.Date.Before(t2))
	}
}

func TestCatRepo_LastUpdates_ShouldReturnLastNApps_Mock(t *testing.T) {
	conn := mockTrackConnection{}
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.LastUpdates(ctx, 123, 4)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 4, len(dboSlice))
}

func TestCatRepo_LastUpdates_ShouldReturnLastNApps(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("app_tracking", "category_tracking")
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	bundleId, _ := testhelpers.AddNewApp(conn, ctx, entities.App{Bundle: "123"})
	track := testhelpers.TrackStruct(bundleId, "key")
	for i := 0; i < 4; i++ {
		_, _ = testhelpers.AddNewTrack(conn, ctx, track, "category_tracking")
	}

	dboSlice, err := repo.LastUpdates(ctx, bundleId, 2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 2, len(dboSlice))

	var key entities.Track
	id := 1000000
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&key))
		assert.Equal(t, "key", key.Type)
		assert.Equal(t, "123", key.App.Bundle)
		assert.Greater(t, id, key.Id)
		id = key.Id
	}
}

func TestCatRepo_ByBundleId_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockTrackConnectionErrors{}
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	_, err := repo.ByBundleId(ctx, 123)
	assert.Error(t, err)
}

func TestCatRepo_ByBundleId_ShouldReturnErrorIfNoRows(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, _ := testhelpers.RealDb(cfg.Database)
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, dboSlice)
}

func TestCatRepo_TimeRange_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockTrackConnectionErrors{}
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	_, err := repo.TimeRange(ctx, 123, time.Now(), time.Now())
	assert.Error(t, err)
}

func TestCatRepo_TimeRange_ShouldReturnErrorIfNoRow(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, _ := testhelpers.RealDb(cfg.Database)
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.TimeRange(ctx, 1, time.Now(), time.Now())
	assert.Error(t, err)
	assert.Nil(t, dboSlice)
}

func TestCatRepo_LastUpdates_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockTrackConnectionErrors{}
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	_, err := repo.LastUpdates(ctx, 123, 1)
	assert.Error(t, err)
}


func TestCatRepo_LastUpdates_ShouldReturnErrorIfNoRow(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, _ := testhelpers.RealDb(cfg.Database)
	repo := trackstore.CatRepo{Conn: conn}
	ctx := context.Background()

	dboSlice, err := repo.LastUpdates(ctx, 1, 2)
	assert.Error(t, err)
	assert.Nil(t, dboSlice)
}



