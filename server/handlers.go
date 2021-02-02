package server

import (
	"Muromachi/auth"
	"Muromachi/graph/generated"
	"Muromachi/httpresp"
	"Muromachi/store"
	"Muromachi/store/entities"
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
func Authorize(sec auth.Defender, sessions *store.AuthCollection) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			request auth.JWTRequest
		)
		// parse body as JWTRequest
		if err := ctx.BodyParser(&request); err != nil {
			return httpresp.Error(ctx, 404, err.Error())
		}
		// CheckAndDel if user with this client id and secret exists
		user, err := sessions.Users.Approve(ctx.Context(), request.ClientId)
		if err != nil {
			return httpresp.Error(ctx, 404, err.Error())
		}
		if err = user.CompareSecret(request.ClientSecret); err != nil {
			return httpresp.Error(ctx, 401, auth.ErrNotAuthenticated)
		}
		// Pass user to request context
		ctx.Locals("request_user", &auth.UserClaims{
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
			return httpresp.Error(ctx, 404, "need to provide access type for request")
		}

		// Create jwt object
		accesstoken, err := sec.SignAccessToken(ctx, refreshToken)
		if err != nil {
			return httpresp.Error(ctx, 400, err.Error())
		}

		// return json depending of the type of Access type
		return ctx.JSON(accesstoken)
	}
}

// Generate new company in system
func NewCompany(sessions *store.AuthCollection) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var err error
		company := ctx.Query("c")
		if company == "" {
			return httpresp.Error(ctx, 404, "empty company query")
		}

		user := entities.User{
			Company:      company,
		}
		err = user.GenerateSecrets()
		if err != nil {
			return httpresp.Error(ctx, 400, "can not generate secrets")
		}
		user, err = sessions.Users.Create(ctx.Context(), user)
		if err != nil {
			return httpresp.Error(ctx, 400, err)
		}

		return ctx.JSON(user)
	}
}

// Ban refresh session
func Ban(collection *store.AuthCollection) func (*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

// Unban refresh session
func Unban(collection *store.AuthCollection) func (*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}