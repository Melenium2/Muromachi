package entities_test

import (
	"Muromachi/graph/model"
	"Muromachi/store/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithAppStruct(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.App{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.App, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithMetaStruct(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Meta{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithKeywordsStruct(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Keywords, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnSliceOfGraphqlFromDBOWithCategoriesStruct(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Categories, len(dboSlice))

	assert.NoError(t, dboSlice.To(modelApps))

	assert.Equal(t, len(dboSlice), len(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfAppIsNotTheSame(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.App{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.App, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfMetaIsNotTheSame(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Meta{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfKeywordsIsNotTheSame(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Keywords, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfLengthOfCategoriesIsNotTheSame(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice) + 2)
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfPassedSliceWithWrongType(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Track{
			Id: i+1,
		})
	}
	modelApps :=  make([]*entities.App, len(dboSlice))
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfDifferentDataTypesInSliceForExampleAppAndTrack(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.App{Id: i+1}, entities.Track{Id: i+1})
	}
	modelApps :=  make([]*model.App, len(dboSlice))
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfDifferentDataTypesInSliceForExampleMetaAndTrack(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Meta{Id: i+1}, entities.Track{Id: i+1})
	}
	modelApps :=  make([]*model.Meta, len(dboSlice))
	assert.Error(t, dboSlice.To(modelApps))
}

func TestDboSlice_To_ShouldReturnErrorIfDifferentDataTypesInSliceForExampleTrackAndApp(t *testing.T) {
	dboSlice := entities.DboSlice{}
	for i := 0; i < 200_000; i++ {
		dboSlice = append(dboSlice, entities.Track{Id: i+1}, entities.App{Id: i+1})
	}
	modelKeys :=  make([]*model.Keywords, len(dboSlice))
	assert.Error(t, dboSlice.To(modelKeys))

	modelCats := make([]*model.Categories, len(dboSlice))
	assert.Error(t, dboSlice.To(modelCats))
}
