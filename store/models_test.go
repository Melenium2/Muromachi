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

	dboMeta := store.Meta{Id: 10, App: store.App{Bundle: "123"}}
	err := dboMeta.To(meta)
	assert.NoError(t, err)

	assert.Equal(t, 10, meta.ID)
	assert.Equal(t, "123", meta.App.Bundle)
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

	dboTrack := store.Track{Id: 10, App: store.App{Bundle: "123"}}
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

	dboApp := store.Track{Id: 10}
	err := dboApp.To(&track)
	assert.Error(t, err)
}

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithAppStruct(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.App{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.App, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithMetaStruct(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Meta{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithKeywordsStruct(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Keywords, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithCategoriesStruct(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Categories, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfAppIsNotTheSame(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.App{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.App, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfMetaIsNotTheSame(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Meta{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfKeywordsIsNotTheSame(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Keywords, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfCategoriesIsNotTheSame(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfPassedSliceWithWrongType(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*store.App, len(dboSlice))
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfDifferentDataTypesInSliceForExampleAppAndTrack(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.App{Id: i+1}, store.Track{Id: i+1})
	}
	modelApps :=  make([]*model.App, len(dboSlice))
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfDifferentDataTypesInSliceForExampleMetaAndTrack(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Meta{Id: i+1}, store.Track{Id: i+1})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice))
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfDifferentDataTypesInSliceForExampleTrackAndApp(t *testing.T) {
	dboSlice := store.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, store.Track{Id: i+1}, store.App{Id: i+1})
	}
	modelKeys :=  make([]*model.Keywords, len(dboSlice))
	assert.Error(t, dboSlice.To(modelKeys))

	modelCats := make([]*model.Categories, len(dboSlice))
	assert.Error(t, dboSlice.To(modelCats))
}

