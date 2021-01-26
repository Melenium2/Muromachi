package store_test

import (
	"Muromachi/graph/model"
	"Muromachi/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDBO_AppShouldGiveAReferenceOfValues(t *testing.T) {
	var app store.App

	dboApp := store.App{Id: 10}
	err := dboApp.To(&app)
	assert.NoError(t, err)

	assert.Equal(t, 10, app.Id)
}

func TestDBO_AppShouldGiveAReferenceOfValuesToGraphqlModel(t *testing.T) {
	app := &model.App{}

	t1 := time.Now()
	dboApp := store.App{Id: 10, StartAt: t1}
	err := dboApp.To(app)
	assert.NoError(t, err)

	assert.Equal(t, 10, app.ID)
	assert.Equal(t, t1, app.StartAt)
}

func TestDBO_App_ShouldReturnErrorIfWrongReference(t *testing.T) {
	var apps []store.App

	dboApp := store.App{Id: 10}
	err := dboApp.To(&apps)
	assert.Error(t, err)
}

func TestDBO_MetaShouldGiveAReferenceOfValues(t *testing.T) {
	var meta store.Meta

	dboMeta := store.Meta{Id: 10}
	err := dboMeta.To(&meta)
	assert.NoError(t, err)

	assert.Equal(t, 10, meta.Id)
}

func TestDBO_MetaShouldGiveAReferenceOfValuesToGraphqlModel(t *testing.T) {
	meta := &model.Meta{}

	dboMeta := store.Meta{Id: 10}
	err := dboMeta.To(meta)
	assert.NoError(t, err)

	assert.Equal(t, 10, meta.ID)
}

func TestDBO_Meta_ShouldReturnErrorIfWrongReference(t *testing.T) {
	var apps []store.Meta

	dboApp := store.Meta{Id: 10}
	err := dboApp.To(&apps)
	assert.Error(t, err)
}

func TestDBO_TrackShouldGiveAReferenceOfValues(t *testing.T) {
	var track store.Track

	dboTrack := store.Track{Id: 10}
	err := dboTrack.To(&track)
	assert.NoError(t, err)

	assert.Equal(t, 10, track.Id)
}

func TestDBO_TrackShouldGiveAReferenceOfValuesToGraphqlModel(t *testing.T) {
	track := &model.Categories{}

	dboTrack := store.Track{Id: 10}
	err := dboTrack.To(track)
	assert.NoError(t, err)

	assert.Equal(t, 10, track.ID)

	newTrack := &model.Keywords{}

	err = dboTrack.To(newTrack)
	assert.NoError(t, err)

	assert.Equal(t, 10, newTrack.ID)
}

func TestDBO_Track_ShouldReturnErrorIfWrongReference(t *testing.T) {
	track := &model.Keywords{}

	dboApp := store.Track{Id: 10}
	err := dboApp.To(&track)
	assert.Error(t, err)
}

