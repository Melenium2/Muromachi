package server

import (
	"Muromachi/auth"
	"Muromachi/config"
	"Muromachi/graph"
	"Muromachi/store"
	"Muromachi/store/connector"
	"Muromachi/store/refreshrepo"
	"Muromachi/store/sessions"
	"Muromachi/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Server struct {
	// Fiber app pointer
	app      *fiber.App
	// Server port for listening
	port     string
	// Application config
	config   config.Config
	// Security interface
	security auth.Defender
	// Graphql resolver
	resolver *graph.Resolver
	// Pointer to auth table collection
	sessions *store.AuthCollection
	// Pointer to tracking tables collection
	tracking *store.TableCollection
}

// Init routes and apply middleware
func (s *Server) initRoutes() {
	//Graphql playground
	s.app.All("/playground", testground())
	// GraphQL Group
	ql := s.app.Group("/ql", auth.ApplyAuthMiddleware(s.security))
	ql.All("/query", graphql(s.resolver))

	// Rest
	// Auth
	s.app.Post("/authorize", Authorize(s.security, s.sessions))
	// Generate new company in system
	urlForGeneration := fmt.Sprintf("/%s/generate", utils.Hash("/generate", 123))
	log.Println("Generation link ", urlForGeneration)
	s.app.Get(urlForGeneration, NewCompany(s.sessions))
}

func (s *Server) Listen() error {
	s.initRoutes()
	return s.app.Listen(s.port)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func New(port string, config config.Config) *Server {
	// Init postgres
	conn, err := connector.EstablishConnection(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	// Pointer to table collection
	tables := store.NewTrackingCollection(conn)

	server := &Server{
		app:    fiber.New(),
		port:   fmt.Sprintf(":%s", port),
		config: config,
		security: auth.NewSecurity(
			config.Auth,
			sessions.New(refreshrepo.New(conn), nil),
		),
		sessions: store.NewAuthCollection(conn),
		tracking: tables,
		resolver: &graph.Resolver{
			Tables: tables,
		},
	}

	return server
}
