package graph

import "github.com/shanmukhsista/go-graphql-starter/pkg/services/notes"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	notesService notes.Service
}

// ProvideNewServerResolver This is the Entry point for our Graphql resolver.
// Any dependencies required by our server must be initialized here.
func ProvideNewServerResolver(notesService notes.Service) *Resolver {
	return &Resolver{notesService: notesService}
}
