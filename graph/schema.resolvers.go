package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"Muromachi/graph/generated"
	"Muromachi/graph/model"
	"context"
	"fmt"
)

func (r *queryResolver) Meta(ctx context.Context) ([]*model.Meta, error) {
	return []*model.Meta {
		{
			Title: "agadaok",
			Description: "PrivetSMarsa",
		},
	}, nil
}

func (r *queryResolver) Cats(ctx context.Context) ([]model.Additional, error) {
	return []model.Additional {
		model.Categories{
			BundleID: 1,
		},
		model.Categories{
			BundleID: 2,
		},
		model.Categories{
			BundleID: 3,
		},
	}, nil
}

func (r *queryResolver) Keys(ctx context.Context) ([]model.Additional, error) {
	return nil, fmt.Errorf("%s", "error ?!")
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
