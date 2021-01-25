package store_test

import (
	"Muromachi/store"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppRepo_ById_ShouldReturnMockApp(t *testing.T) {
	conn := mockAppConnection{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	dbo, err := repo.ById(ctx, 10)
	assert.NoError(t, err)
	assert.NotNil(t, dbo)
	var app store.App
	assert.NoError(t, dbo.To(&app))
	assert.Equal(t, "FINANCE", app.Category)
}

func TestAppRepo_ByBundleId_ShouldReturnSliceOfMockApps(t *testing.T) {
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

func TestAppRepo_ById_ShouldReturnErrNoRows(t *testing.T) {
	conn := mockAppConnectionErrors{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	_, err := repo.ById(ctx, 10)
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
}

func TestAppRepo_ByBundleId_ShouldReturnErrNoRows(t *testing.T) {
	conn := mockAppConnectionErrors{}
	repo := store.NewApp(conn)
	ctx := context.Background()

	_, err := repo.ByBundleId(ctx, 10)
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
}
