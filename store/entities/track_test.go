package entities_test

import (
	"Muromachi/graph/model"
	"Muromachi/store/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBO_TrackShouldGiveAReferenceOfValues(t *testing.T) {
	var track entities.Track

	dboTrack := entities.Track{Id: 10}
	err := dboTrack.To(&track)
	assert.NoError(t, err)

	assert.Equal(t, 10, track.Id)
}

func TestDBO_TrackShouldGiveAReferenceOfValuesToGraphqlModel(t *testing.T) {
	track := &model.Categories{}

	dboTrack := entities.Track{Id: 10, App: entities.App{Bundle: "123"}}
	err := dboTrack.To(track)
	assert.NoError(t, err)

	assert.Equal(t, 10, track.ID)
	assert.Equal(t, "123", track.App.Bundle)

	newTrack := &model.Keywords{}

	err = dboTrack.To(newTrack)
	assert.NoError(t, err)

	assert.Equal(t, 10, newTrack.ID)
	assert.Equal(t, "123", newTrack.App.Bundle)
}

func TestDBO_Track_ShouldReturnErrorIfWrongReference(t *testing.T) {
	track := &model.Keywords{}

	dboApp := entities.Track{Id: 10}
	err := dboApp.To(&track)
	assert.Error(t, err)
}
