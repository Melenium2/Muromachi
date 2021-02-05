package server

import (
	"Muromachi/auth"
	"Muromachi/config"
	"Muromachi/graph"
	"Muromachi/store/connector"
	tracking2 "Muromachi/store/tracking"
	"Muromachi/store/users"
	"Muromachi/store/users/sessions"
	"Muromachi/store/users/sessions/blacklist"
	"Muromachi/store/users/sessions/tokens"
	"Muromachi/store/users/userstore"
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
	sessions *users.Tables
	// Pointer to tracking tables collection
	tracking *tracking2.Tables
}

// Init routes and apply middleware
func (s *Server) initRoutes() {
	//Graphql playground
	s.app.All("/playground", Testground())
	// GraphQL Group
	ql := s.app.Group("/ql", auth.ApplyAuthMiddleware(s.security))
	ql.All("/query", Graphql(s.resolver))

	// Rest
	// Auth
	s.app.Post("/authorize", Authorize(s.security, s.sessions))
	// Generate new company in system
	urlForGeneration := fmt.Sprintf("/%s/generate", utils.Hash("/generate", 123))
	log.Println("Generation link ", urlForGeneration)
	s.app.Get(urlForGeneration, NewCompany(s.sessions))
}

// Start listening tcp port
func (s *Server) Listen() error {
	s.initRoutes()
	return s.app.Listen(s.port)
}

// Shutdown server
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// Init new server with given port and application config
func New(port string, config config.Config) *Server {
	// Init postgres
	conn, err := connector.EstablishPostgresConnection(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	// Init redis
	redisConn := connector.EstablishRedisConnection(config.Database.Redis)
	// Pointer to table collection
	tables := tracking2.NewTrackingTables(conn)
	// Interface of sessions
	session := sessions.New(tokens.New(conn), blacklist.New(redisConn))

	server := &Server{
		app:    fiber.New(),
		port:   fmt.Sprintf(":%s", port),
		config: config,
		security: auth.NewSecurity(
			config.Auth,
			session,
		),
		sessions: users.NewAuthTables(session, userstore.NewUserRepo(conn)),
		tracking: tables,
		resolver: &graph.Resolver{
			Tables: tables,
		},
	}

	return server
}
