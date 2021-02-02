package entities_test

import (
	"Muromachi/graph/model"
	"Muromachi/store/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDBO_MetaShouldGiveAReferenceOfValues(t *testing.T) {
	var meta entities.Meta

	dboMeta := entities.Meta{Id: 10}
	err := dboMeta.To(&meta)
	assert.NoError(t, err)

	assert.Equal(t, 10, meta.Id)
}

func TestDBO_MetaShouldGiveAReferenceOfValuesToGraphqlModel(t *testing.T) {
	meta := &model.Meta{}

	dboMeta := entities.Meta{Id: 10, App: entities.App{Bundle: "123"}}
	err := dboMeta.To(meta)
	assert.NoError(t, err)

	assert.Equal(t, 10, meta.ID)
	assert.Equal(t, "123", meta.App.Bundle)
}

func TestDBO_Meta_ShouldReturnErrorIfWrongReference(t *testing.T) {
	var apps []entities.Meta

	dboApp := entities.Meta{Id: 10}
	err := dboApp.To(&apps)
	assert.Error(t, err)
}

