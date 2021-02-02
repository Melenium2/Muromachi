package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"Muromachi/graph/generated"
	"Muromachi/graph/model"
	"Muromachi/graph/scalar"
	"Muromachi/store/entities"
	"context"
	"time"
)

func (r *queryResolver) Meta(ctx context.Context, id int, last *int, start *scalar.FormattedDate, end *scalar.FormattedDate) ([]*model.Meta, error) {
	var (
		dbo entities.DboSlice
		err error
	)
	if start != nil && end != nil {
		dbo, err = r.Tables.Meta.TimeRange(ctx, id, time.Time(*start), time.Time(*end))
		if err != nil {
			return nil, err
		}
	} else if last != nil {
		dbo, err = r.Tables.Meta.LastUpdates(ctx, id, *last)
		if err != nil {
			return nil, err
		}
	} else {
		dbo, err = r.Tables.Meta.ByBundleId(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	metaModels := make([]*model.Meta, len(dbo))
	if err := dbo.To(metaModels); err != nil {
		return nil, err
	}

	return metaModels, nil
}

func (r *queryResolver) Cats(ctx context.Context, id int, last *int, start *scalar.FormattedDate, end *scalar.FormattedDate) ([]*model.Categories, error) {
	var (
		dbo entities.DboSlice
		err error
	)
	if start != nil && end != nil {
		dbo, err = r.Tables.Cat.TimeRange(ctx, id, time.Time(*start), time.Time(*end))
		if err != nil {
			return nil, err
		}
	} else if last != nil {
		dbo, err = r.Tables.Cat.LastUpdates(ctx, id, *last)
		if err != nil {
			return nil, err
		}
	} else {
		dbo, err = r.Tables.Cat.ByBundleId(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	metaModels := make([]*model.Categories, len(dbo))
	if err := dbo.To(metaModels); err != nil {
		return nil, err
	}

	return metaModels, nil
}

func (r *queryResolver) Keys(ctx context.Context, id int, last *int, start *scalar.FormattedDate, end *scalar.FormattedDate) ([]*model.Keywords, error) {
	var (
		dbo entities.DboSlice
		err error
	)
	if start != nil && end != nil {
		dbo, err = r.Tables.Keys.TimeRange(ctx, id, time.Time(*start), time.Time(*end))
		if err != nil {
			return nil, err
		}
	} else if last != nil {
		dbo, err = r.Tables.Keys.LastUpdates(ctx, id, *last)
		if err != nil {
			return nil, err
		}
	} else {
		dbo, err = r.Tables.Keys.ByBundleId(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	metaModels := make([]*model.Keywords, len(dbo))
	if err := dbo.To(metaModels); err != nil {
		return nil, err
	}

	return metaModels, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
