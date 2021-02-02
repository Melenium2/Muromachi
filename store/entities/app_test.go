package entities_test

import (
	"Muromachi/graph/model"
	"Muromachi/store/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDBO_AppShouldGiveAReferenceOfValues(t *testing.T) {
	var app entities.App

	dboApp := entities.App{Id: 10}
	err := dboApp.To(&app)
	assert.NoError(t, err)

	assert.Equal(t, 10, app.Id)
}

func TestDBO_AppShouldGiveAReferenceOfValuesToGraphqlModel(t *testing.T) {
	app := &model.App{}

	t1 := time.Now()
	dboApp := entities.App{Id: 10, StartAt: t1}
	err := dboApp.To(app)
	assert.NoError(t, err)

	assert.Equal(t, 10, app.ID)
	assert.Equal(t, t1, app.StartAt)
}

func TestDBO_App_ShouldReturnErrorIfWrongReference(t *testing.T) {
	var apps []entities.App

	dboApp := entities.App{Id: 10}
	err := dboApp.To(&apps)
	assert.Error(t, err)
}
