package authorization

import (
	"Muromachi/server"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Authentication middleware for service
func ApplyAuthMiddleware(security Defender) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var token string
		// Check if token  in cookie
		cookieToken := c.Cookies(SecurityCookieName, "")
		if cookieToken == "" {
			// if not in cookie then check headers
			authToken := c.Get("Authorization", "")
			if authToken == "" {
				// cookie and headers empty -> return err
				return server.HttpError(c, 401, ErrNotAuthenticated)
			} else {
				token = authToken[7:]
			}
		} else {
			token = cookieToken
		}

		// Validate jwt
		claims, err := security.ValidateJwt(token)
		if err != nil {
			ErrNotAuthenticated["additional"] = err
			return server.HttpError(c, 401, ErrNotAuthenticated)
		}

		// Check if refresh token is banned in redis
		// TODO Сделать Redis
		if security.IsSessionBanned(claims.Id) {
			return server.HttpError(c, 401, ErrNotAuthenticated)
		}

		expiredTime := time.Unix(claims.ExpiresAt, 0)
		// Is Jwt token expired ?
		if expiredTime.After(time.Now()) {
			ErrNotAuthenticated["additional"] = ErrExpiredAccessToken
			return server.HttpError(c, 401, ErrNotAuthenticated)
		}

		c.Locals("request_user", claims.UserClaims)

		return c.Next()
	}
}
