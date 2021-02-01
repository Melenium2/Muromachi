package server

import (
	"Muromachi/authorization"
	"Muromachi/graph/generated"
	"Muromachi/httpresp"
	"Muromachi/store"
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
func authorize(sec authorization.Defender, sessions *store.AuthCollection) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			request     authorization.JWTRequest
		)
		// parse body as JWTRequest
		if err := ctx.BodyParser(&request); err != nil {
			return httpresp.Error(ctx, 404, err)
		}
		// Check if userrepo with this client id and secret exists
		user, err := sessions.Users.Approve(ctx.Context(), request.ClientId)
		if err != nil {
			return httpresp.Error(ctx, 404, err)
		}
		if err = user.CompareSecret(request.ClientSecret); err != nil {
			return httpresp.Error(ctx, 401, authorization.ErrNotAuthenticated)
		}
		// Pass userrepo to request context
		ctx.Locals("request_user", &authorization.UserClaims{
			ID:   int64(user.ID),
			Role: "user",
		})

		// depending of the type chose response params
		// if type session or refresh_token then get new refresh token
		var refreshToken string
		switch request.AccessType {
		case "simple":
			break
		case "session":
			refreshToken, err = sec.StartSession(ctx)
			if err != nil {
				return httpresp.Error(ctx, 400, err)
			}
		case "refresh_token":
			refreshToken, err = sec.StartSession(ctx, request.RefreshToken)
			if err != nil {
				return httpresp.Error(ctx, 400, err)
			}
		default:
			return httpresp.Error(ctx, 404, map[string]interface{}{
				"error": "need to provide access type for request",
			})
		}

		// Create jwt object
		accesstoken, err := sec.SignAccessToken(ctx, refreshToken)
		if err != nil {
			return httpresp.Error(ctx, 400, err)
		}

		// return json depending of the type of Access type
		return ctx.JSON(accesstoken)
	}
}

