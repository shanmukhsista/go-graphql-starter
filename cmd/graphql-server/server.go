package main

import (
	"fmt"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/shanmukhsista/go-graphql-starter/cmd/graphql-server/dependencies"
	"github.com/shanmukhsista/go-graphql-starter/cmd/graphql-server/graph"
	"github.com/shanmukhsista/go-graphql-starter/cmd/graphql-server/graph/generated"
	"github.com/shanmukhsista/go-graphql-starter/pkg/common/config"
)

const defaultPort = "8080"

func main() {
	port := GetPortFromEnv()
	// get config file from go args -configpath
	configPath := config.MustGetConfigPathFromFlags("configpath")
	srv, err := SetupServer(configPath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to start server")
		return
	}
	err = srv.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to start server on port %s", port)
	}
}

const allowedOriginsConfigKey = "server.cors.allowed_origins"

// Defining the Graphql handler
func NewGraphqlHandler(resolver *graph.Resolver) echo.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := GetConfigForGrqphQLServer(resolver)
	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func GetPortFromEnv() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}

func GetConfigForGrqphQLServer(resolver *graph.Resolver) generated.Config {
	c := generated.Config{Resolvers: resolver}
	return c
}

// Defining the Playground handler
func playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func setCORSMiddleware(r *echo.Echo) {
	allowedOrigins := config.MustGetStringSet(allowedOriginsConfigKey)
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{"PUT", "POST", "GET", "DELETE"},
		AllowHeaders: []string{"*"},
		MaxAge:       int((12 * time.Hour).Seconds()),
		AllowOriginFunc: func(origin string) (bool, error) {
			return lo.Contains(allowedOrigins, origin), nil
		},
	}))

}

func SetupServer(configPath string) (*echo.Echo, error) {
	err := config.MustLoadConfigAtPath(configPath)
	r := echo.New()

	if err != nil {
		return nil, err
	}
	resolver, err := dependencies.NewAppResolverService()
	if err != nil {
		return nil, err
	}
	// Add cors middleware.
	setCORSMiddleware(r)

	r.POST("/query", NewGraphqlHandler(resolver))
	// TODO only in dev mode.
	r.GET("/", playgroundHandler())

	return r, nil
}
