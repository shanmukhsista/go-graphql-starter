//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/google/wire"
	"github.com/shanmukhsista/go-graphql-starter/cmd/graphql-server/graph"
	"github.com/shanmukhsista/go-graphql-starter/pkg/common/db"
	"github.com/shanmukhsista/go-graphql-starter/pkg/services/notes"
)

var postgresDbConnectionSet = wire.NewSet(db.ProvidePgConnectionPool,
	db.ProvideNewPostgresTransactor, db.ProvideNewDatabaseConnection)

var notesApiService = wire.NewSet(notes.ProvideNewNotesRepository, notes.ProvideNewNotesService)

var graphQLServerDependencySet = wire.NewSet(
	postgresDbConnectionSet,
	notesApiService,
	graph.ProvideNewServerResolver,
)

func NewAppResolverService() (*graph.Resolver, error) {
	wire.Build(graphQLServerDependencySet)
	return nil, nil
}
