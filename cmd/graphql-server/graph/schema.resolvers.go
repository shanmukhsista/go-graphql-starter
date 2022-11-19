package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/shanmukhsista/go-graphql-starter/cmd/graphql-server/graph/generated"
	"github.com/shanmukhsista/go-graphql-starter/pkg/model"
)

// CreateNewNote is the resolver for the createNewNote field.
func (r *mutationResolver) CreateNewNote(ctx context.Context, input model.NewNoteInput) (*model.Note, error) {
	return r.notesService.SaveNewNote(ctx, input)
}

// Notes is the resolver for the notes field.
func (r *queryResolver) Notes(ctx context.Context) ([]*model.Note, error) {
	return r.notesService.GetAllNotes(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
