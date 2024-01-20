package main

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/tkdn/go-study/domain"
	"github.com/tkdn/go-study/graph"
	"github.com/tkdn/go-study/infra"
	"github.com/tkdn/go-study/log"
	"github.com/tkdn/go-study/telemetry"
	"github.com/tkdn/go-study/web/middleware"
)

func main() {
	ctx := context.Background()
	db, err := infra.ConnectDB()
	if err != nil {
		log.Logger.Error(err.Error())
	}
	defer db.Close()

	tpShutdown, err := telemetry.Do(ctx)
	if err != nil {
		log.Logger.Error("failed to start telemetry.")
	}
	defer tpShutdown()

	c := graph.Config{
		Resolvers: &graph.Resolver{
			UserRepo: domain.NewUserRepository(db),
			PostRepo: domain.NewPostRepository(db),
		},
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))
	mux := http.NewServeMux()
	mux.Handle("/playground", playground.Handler("Graphql playground", "/graphql"))
	mux.Handle("/graphql", srv)
	otelm := telemetry.NewOtelHttpMiddleware()

	server := &http.Server{
		Addr:    net.JoinHostPort("", "8080"),
		Handler: otelm(middleware.RejectAdminInternal(mux)),
	}
	log.Logger.Info("server start")
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Logger.Error(err.Error())
	}
}
