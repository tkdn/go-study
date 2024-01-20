package main

import (
	"errors"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/tkdn/go-study/domain"
	"github.com/tkdn/go-study/graph"
	"github.com/tkdn/go-study/infra"
	"github.com/tkdn/go-study/log"
	"github.com/tkdn/go-study/middleware"
)

func main() {
	db := infra.ConnectDB()
	defer db.Close()

	c := graph.Config{
		Resolvers: &graph.Resolver{
			UserRepo: domain.NewUserRepository(db),
		},
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))
	mux := http.NewServeMux()
	mux.Handle("/playground", playground.Handler("Graphql playground", "/graphql"))
	mux.Handle("/graphql", srv)

	server := &http.Server{
		Addr:    net.JoinHostPort("", "8080"),
		Handler: middleware.RejectAdminInternal(mux),
	}
	log.Logger.Info("server start")
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Logger.Error(err.Error())
	}
}
