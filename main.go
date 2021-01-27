package main

import (
	"Muromachi/config"
	"Muromachi/graph"
	"Muromachi/graph/generated"
	"Muromachi/store"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"os"
)

// Даталоадер, агрегация постоянных одинаковых запросов
// https://gqlgen.com/reference/dataloaders/
//
// Сделать авторизацию (обсудить как удобнее) 	https://graphql.org/learn/authorization/
//												https://gqlgen.com/recipes/authentication/
// Показывать или нет ендпоинты https://gqlgen.com/reference/introspection/

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cfg := config.New("./config/dev.yml")
	cfg.Database.Schema = "./config/schema.sql"
	conn, err := store.EstablishConnection(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	resolver := &graph.Resolver{
		Tables: store.New(conn),
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
