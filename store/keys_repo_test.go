package store_test

import (
	"Muromachi/store"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestKeysRepo_ByBundleId_ShouldReturnApp_Mock(t *testing.T) {
	conn := mockTrackConnection{}
	repo := store.NewKeys(conn)
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 123)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var key store.Track
	assert.NoError(t, dboSlice[0].To(&key))
	assert.Equal(t, "type", key.Type)
}

func TestKeysRepo_ByBundleId_ShouldReturnApp(t *testing.T) {
	conn, cleaner := RealDb()
	defer cleaner("app_tracking", "keyword_tracking")
	repo := store.NewKeys(conn)
	ctx := context.Background()

	bundleId, _ := AddNewApp(conn, ctx, store.App{Bundle: "123"})
	track := TrackStruct(bundleId, "key")
	for i := 0; i < 4; i++ {
		_, _ = AddNewTrack(conn, ctx, track, "keyword_tracking")
	}

	dboSlice, err := repo.ByBundleId(ctx, bundleId)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 4, len(dboSlice))

	var key store.Track
	assert.NoError(t, dboSlice[0].To(&key))
	assert.Equal(t, "key", key.Type)
}

func TestKeysRepo_TimeRange_ShouldReturnAppsInTimeRange_Mock(t *testing.T) {
	conn := mockTrackConnection{}
	repo := store.NewKeys(conn)
	ctx := context.Background()

	t1, _ := time.Parse("2006-01-02", "2021-01-18")
	t2 := t1.AddDate(0, 0, 10)
	dboSlice, err := repo.TimeRange(ctx, 123, t1, t2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)

	var key store.Track
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&key))
		assert.Equal(t, "type", key.Type)
		assert.True(t, key.Date.After(t1) && key.Date.Before(t2))
	}
}

func TestKeysRepo_TimeRange_ShouldReturnAppsInTimeRange(t *testing.T) {
	conn, cleaner := RealDb()
	defer cleaner("app_tracking", "keyword_tracking")
	repo := store.NewKeys(conn)
	ctx := context.Background()

	bundleId, _ := AddNewApp(conn, ctx, store.App{Bundle: "123"})
	track := TrackStruct(bundleId, "key")
	t1 := track.Date.AddDate(0, 0, -1)
	for i := 0; i < 4; i++ {
		_, _ = AddNewTrack(conn, ctx, track, "keyword_tracking")
		track.Date = track.Date.AddDate(0, 0, 1)
	}
	t2 := track.Date.AddDate(0, 0, 1)

	dboSlice, err := repo.TimeRange(ctx, bundleId, t1, t2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 4, len(dboSlice))

	var key store.Track
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&key))
		assert.Equal(t, "key", key.Type)
		assert.True(t, key.Date.After(t1) && key.Date.Before(t2))
	}
}

func TestKeysRepo_LastUpdates_ShouldReturnLastNApps_Mock(t *testing.T) {
	conn := mockTrackConnection{}
	repo := store.NewKeys(conn)
	ctx := context.Background()

	dboSlice, err := repo.LastUpdates(ctx, 123, 4)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 4, len(dboSlice))
}

func TestKeysRepo_LastUpdates_ShouldReturnLastNApps(t *testing.T) {
	conn, cleaner := RealDb()
	defer cleaner("app_tracking", "keyword_tracking")
	repo := store.NewKeys(conn)
	ctx := context.Background()

	bundleId, _ := AddNewApp(conn, ctx, store.App{Bundle: "123"})
	track := TrackStruct(bundleId, "key")
	for i := 0; i < 4; i++ {
		_, _ = AddNewTrack(conn, ctx, track, "keyword_tracking")
	}

	dboSlice, err := repo.LastUpdates(ctx, bundleId, 2)
	assert.NoError(t, err)
	assert.NotNil(t, dboSlice)
	assert.Equal(t, 2, len(dboSlice))

	var key store.Track
	id := 1000
	for _, v := range dboSlice {
		assert.NoError(t, v.To(&key))
		assert.Equal(t, "key", key.Type)
		assert.Greater(t, id, key.Id)
		id = key.Id
	}
}

func TestKeysRepo_ByBundleId_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockTrackConnectionErrors{}
	repo := store.NewKeys(conn)
	ctx := context.Background()

	_, err := repo.ByBundleId(ctx, 123)
	assert.Error(t, err)
}

func TestKeysRepo_ByBundleId_ShouldReturnErrorIfNoRows(t *testing.T) {
	conn, _ := RealDb()
	repo := store.NewKeys(conn)
	ctx := context.Background()

	dboSlice, err := repo.ByBundleId(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, dboSlice)
}

func TestKeysRepo_TimeRange_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockTrackConnectionErrors{}
	repo := store.NewKeys(conn)
	ctx := context.Background()

	_, err := repo.TimeRange(ctx, 123, time.Now(), time.Now())
	assert.Error(t, err)
}

func TestKeysRepo_TimeRange_ShouldReturnErrorIfNoRow(t *testing.T) {
	conn, _ := RealDb()
	repo := store.NewKeys(conn)
	ctx := context.Background()

	dboSlice, err := repo.TimeRange(ctx, 1, time.Now(), time.Now())
	assert.Error(t, err)
	assert.Nil(t, dboSlice)
}

func TestKeysRepo_LastUpdates_ShouldReturnErrorIfNoRows_Mock(t *testing.T) {
	conn := mockTrackConnectionErrors{}
	repo := store.NewKeys(conn)
	ctx := context.Background()

	_, err := repo.LastUpdates(ctx, 123, 1)
	assert.Error(t, err)
}


func TestKeysRepo_LastUpdates_ShouldReturnErrorIfNoRow(t *testing.T) {
	conn, _ := RealDb()
	repo := store.NewKeys(conn)
	ctx := context.Background()

	dboSlice, err := repo.LastUpdates(ctx, 1, 2)
	assert.Error(t, err)
	assert.Nil(t, dboSlice)
}


