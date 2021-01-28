package server

import (
	"Muromachi/authorization"
	"Muromachi/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
)

func testground() func(ctx *fiber.Ctx) error {
	play := playground.Handler("GraphQL playground", "/query")

	return func(ctx *fiber.Ctx) error {
		play(ctx.Context())
		return nil
	}
}

func graphql(resolver generated.ResolverRoot) func(ctx *fiber.Ctx) error {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: resolver},
		),
	).Handler()

	return func(ctx *fiber.Ctx) error {
		srv(ctx.Context())
		return nil
	}
}

// Auth endpoint
func (s *Server) authorize(ctx *fiber.Ctx) error {
	var (
		request     authorization.JWTRequest
	)
	// parse body as JWTRequest
	if err := ctx.BodyParser(&request); err != nil {
		ctx.Status(404)
		return err
	}
	// Check if user with this client id and secret exists
	user, err := s.sessions.Users.Approve(ctx.Context(), request.ClientId, request.ClientSecret)
	if err != nil {
		ctx.Status(404)
		return err
	}
	// Pass user to request context
	ctx.Locals("request_user", &authorization.UserClaims{
		ID:   int64(user.ID),
		Role: "user",
	})

	// Check type if request
	var refreshToken string

	// TODO Разобраться как правильно поступить в дданной ситуации
	switch request.AccessType {
	case "simple":
		break
	case "session":
		refreshToken, err = s.security.StartSession(ctx)
		if err != nil {
			return err
		}
	case "refresh_token":
		refreshToken, err = s.security.StartSession(ctx, request.RefreshToken)
		if err != nil {
			return err
		}
	default:
		ctx.Status(404)
		return ctx.JSON(map[string]interface{}{
			"error": "need to provide access type for request",
		})
	}

	// Create jwt object
	accesstoken, err := s.security.SignAccessToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return ctx.JSON(accesstoken)
}