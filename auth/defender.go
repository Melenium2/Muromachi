package auth

import (
	"Muromachi/config"
	"Muromachi/store/entities"
	"Muromachi/store/users/sessions"
	"context"
	"github.com/gofiber/fiber/v2"
)

// Defender is interface for session management
type Defender interface {
	StartSession(ctx *fiber.Ctx, refreshToken ...string) (string, error)
	IsSessionBanned(ctx context.Context, refreshToken string) bool
	BanSessions(ctx context.Context, tokens ...entities.Session) error
	SignAccessToken(ctx *fiber.Ctx, refreshToken string) (JWTResponse, error)
	ValidateJwt(accessToken string) (*Claims, error)
}

func NewSecurity(config config.Authorization, usersession sessions.Session) *Security {
	return &Security{
		config:    config,
		generator: newSecurityGen(config),
		sessions:  usersession,
	}
}
