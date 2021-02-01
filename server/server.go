package server

import (
	"Muromachi/authorization"
	"Muromachi/config"
	"Muromachi/graph"
	"Muromachi/store"
	"Muromachi/store/connector"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Server struct {
	app      *fiber.App
	port     string
	config   config.Config
	security *authorization.Security
	resolver *graph.Resolver
	sessions *store.AuthCollection
	tracking *store.TableCollection
}

// Init routes and apply middleware
func (s *Server) InitRoutes() {
	//Graphql playground
	s.app.All("/playground", testground())
	// GraphQL Group
	ql := s.app.Group("/ql", s.security.ApplyAuthMiddleware)
	ql.All("/query", graphql(s.resolver))

	// Rest
	// Auth
	s.app.Post("/authorize", s.authorize)
}

func (s *Server) Listen() error {
	s.InitRoutes()
	return s.app.Listen(s.port)
}

func New(port string, config config.Config) *Server {
	// Init postgres
	config.Database.Schema = "../config/schema.sql"
	conn, err := connector.EstablishConnection(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	// Pointer to table collection
	tables := store.New(conn)

	return &Server{
		app:      fiber.New(),
		port:     fmt.Sprintf(":%s", port),
		config:   config,
		security: nil,
		tracking: tables,
		resolver: &graph.Resolver{
			Tables: tables,
		},
	}
}
