package server

import (
	"Muromachi/auth"
	"Muromachi/graph/generated"
	"Muromachi/httpresp"
	"Muromachi/server/requests"
	"Muromachi/store/entities"
	"Muromachi/store/users"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Handler for graphql testground
func Testground() func(ctx *fiber.Ctx) error {
	play := playground.Handler("GraphQL playground", "/query")

	return func(ctx *fiber.Ctx) error {
		play(ctx.Context())
		return nil
	}
}

// Graphql handler
func Graphql(resolver generated.ResolverRoot) func(ctx *fiber.Ctx) error {
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
func Authorize(sec auth.Defender, sessions *users.Tables) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			request auth.JWTRequest
		)
		// parse body as JWTRequest
		if err := ctx.BodyParser(&request); err != nil {
			return httpresp.Error(ctx, 404, err.Error())
		}
		if request.ClientId == "" || request.ClientSecret == "" {
			return httpresp.Error(ctx, 404, fmt.Errorf("%s", "client id or client secret not provided"))
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
func NewCompany(sessions *users.Tables) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var err error
		company := ctx.Query("c")
		if company == "" {
			return httpresp.Error(ctx, 404, "empty company query")
		}

		user := entities.User{
			Company: company,
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
func Ban(collection *users.Tables) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			list   requests.TokenList
			forBan []entities.Session
		)
		if err := ctx.BodyParser(&list); err != nil {
			return httpresp.Error(ctx, 400, "can not parse to token list")
		}

		// If user id is presented
		if list.UserId > 0 {
			sessions, err := collection.Sessions.UserSessions(ctx.Context(), list.UserId)
			if err == nil {
				forBan = append(forBan, sessions...)
			}
		}

		// Adding sessions with same tokens from token list
		for _, token := range list.Tokens {
			session, err := collection.Sessions.Get(ctx.Context(), token)
			if err != nil {
				return httpresp.Error(ctx, 400, err.Error())
			}
			forBan = append(forBan, session)
		}

		// Add sessions to blacklist
		for _, s := range forBan {
			if err := collection.Sessions.Add(
				ctx.Context(),
				s.RefreshToken,
				s.ID,
				time.Duration(list.Ttl),
			); err != nil {
				return httpresp.Error(ctx, 400, err.Error())
			}
		}

		return ctx.JSON(requests.BanInfo{
			Type:   "ban",
			Count:  len(forBan),
			Tokens: forBan,
			At:     time.Now(),
		})
	}
}

// Unban refresh session
func Unban(collection *users.Tables) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			list    requests.TokenList
			antiBan []string
		)
		if err := ctx.BodyParser(&list); err != nil {
			return httpresp.Error(ctx, 400, "can not parse to token list")
		}

		// if user id is presented
		if list.UserId > 0 {
			sessions, err := collection.Sessions.UserSessions(ctx.Context(), list.UserId)
			if err == nil {
				for _, s := range sessions {
					antiBan = append(antiBan, s.RefreshToken)
				}
			}
		}
		// append given tokens from request
		antiBan = append(antiBan, list.Tokens...)

		// if slice not nil then remove sessions from blacklist
		if len(antiBan) > 0 {
			n, err := collection.Sessions.Del(ctx.Context(), antiBan...)
			if err != nil || int(n) != len(antiBan) {
				return httpresp.Error(ctx, 500, "unexpected error while deletion")
			}
		}

		return ctx.JSON(requests.BanInfo{
			Type:   "unban",
			Count:  len(antiBan),
			Tokens: antiBan,
			At:     time.Now(),
		})
	}
}
