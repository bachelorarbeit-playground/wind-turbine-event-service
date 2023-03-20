package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"context"
	"encoding/json"
	"gql-service/graph/model"
	"gql-service/pkg/storage"

	"github.com/rs/zerolog/log"
)

// AverageProduction is the resolver for the AverageProduction field.
func (r *queryResolver) AverageProduction(ctx context.Context) ([]*model.AverageProductionEntry, error) {
	return getEntries[model.AverageProductionEntry](storage.Store.AverageProduction), nil
}

// AnomalyDetection is the resolver for the AnomalyDetection field.
func (r *queryResolver) AnomalyDetection(ctx context.Context) ([]*model.AnomalyDetectionEntry, error) {
	return getEntries[model.AnomalyDetectionEntry](storage.Store.AnomalyDetection), nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func getEntries[S any, T any](view map[string]*T) []*S {
	entries := make([]*S, 0, len(view))

	for _, value := range view {
		bytes, err := json.Marshal(value)
		if err != nil {
			log.Panic().Err(err)
		}

		var result S
		if err = json.Unmarshal(bytes, &result); err != nil {
			log.Panic().Err(err)
		}

		entries = append(entries, &result)
	}

	return entries
}
