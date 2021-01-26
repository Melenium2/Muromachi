package main

import (
	"Muromachi/graph"
	"Muromachi/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"os"
)

// Сделать кеширование https://graphql.org/learn/caching/
// Сделать авторизацию (обсудить как удобнее) 	https://graphql.org/learn/authorization/
//												https://gqlgen.com/recipes/authentication/
// Сделать error log https://gqlgen.com/reference/errors/
// Показывать или нет ендпоинты https://gqlgen.com/reference/introspection/
// Поставить ограничение по сложности запросов https://gqlgen.com/reference/complexity/

// База данных
// Генерация бд при условии что ее нет
// Создать индексы в бд
// Дописать тесты с реальной базой для APP_REPO

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
