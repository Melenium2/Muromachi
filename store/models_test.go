package store_test

import (
	"Muromachi/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBO_AppShouldGiveAReferenceOfValues(t *testing.T) {
	var app store.App

	dboApp := store.App{Id: 10}
	err := dboApp.To(&app)
	assert.NoError(t, err)

	assert.Equal(t, 10, app.Id)
}

func TestDBO_MetaShouldGiveAReferenceOfValues(t *testing.T) {
	var meta store.Meta

	dboMeta := store.Meta{Id: 10}
	err := dboMeta.To(&meta)
	assert.NoError(t, err)

	assert.Equal(t, 10, meta.Id)
}

func TestDBO_TrackShouldGiveAReferenceOfValues(t *testing.T) {
	var track store.Track

	dboTrack := store.Track{Id: 10}
	err := dboTrack.To(&track)
	assert.NoError(t, err)

	assert.Equal(t, 10, track.Id)
}