package auth

import (
	"Muromachi/config"
	"Muromachi/store/entities"
	"Muromachi/store/sessions"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"time"
)

var (
	SecurityCookieName = "apptwice-access-token"

	ErrNotAuthenticated = map[string]interface{}{
		"status": 401,
		"error":  "invalid auth token, please login with you credentials",
	}
)

type Defender interface {
	StartSession(ctx *fiber.Ctx, refreshToken ...string) (string, error)
	IsSessionBanned(ctx context.Context, refreshToken string) bool
	BanSessions(ctx context.Context, tokens ...entities.Session) error
	SignAccessToken(ctx *fiber.Ctx, refreshToken string) (JWTResponse, error)
	ValidateJwt(accessToken string) (*Claims, error)
}

type Security struct {
	config    config.Authorization
	generator *securityGenerator
	sessions  sessions.Sessions
}

func (security *Security) ApplyRequestIdMiddleware(c *fiber.Ctx) error {
	// TODO Добавить реквест id к каждому запросу. Переместить метод в другой файл
	return c.Next()
}

// Creating new refresh session in DB and return new refresh token for user
//
// TODO Возможно стоит сравнивать user-agent или дополнительные парамтры для ваидации
func (security *Security) StartSession(ctx *fiber.Ctx, refreshToken ...string) (string, error) {
	var (
		token  string
		userId int
	)
	// If we recreating existed refresh session
	if len(refreshToken) > 0 && refreshToken[0] != "" {
		token = refreshToken[0]
		// Remove old session and get her value
		session, err := security.sessions.Remove(ctx.Context(), token)
		if err != nil {
			if err == pgx.ErrNoRows {
				return "", fmt.Errorf("%s", "session not found")
			}
			return "", err
		}
		// CheckAndDel if session not expired
		if time.Now().After(session.ExpiresIn) {
			return "", ErrExpiredRefreshToken
		}
		if security.IsSessionBanned(ctx.Context(), token) {
			return "", fmt.Errorf("%s", "your refresh session in blacklist")
		}
		// Save userid for creating new session
		userId = session.UserId
	}

	// If we approve user with his credentials, we have *UserClaims stored inside ctx
	if userId == 0 {
		claims, ok := ctx.Locals("request_user").(*UserClaims)
		if !ok {
			return "", ErrEmptyContext
		}
		userId = int(claims.ID)
	}
	// CheckAndDel if user has < 5 opened sessions
	userSessions, err := security.sessions.UserSessions(ctx.Context(), userId)
	if len(userSessions) > 5 {
		// Map only ids of sessions
		s := make([]int, len(userSessions))
		for i, v := range userSessions {
			s[i] = v.ID
		}
		// Then remove all this sessions
		if err = security.sessions.RemoveBatch(ctx.Context(), s...); err != nil {
			return "", err
		}
	}

	// Create new session in DB
	newSession, err := security.sessions.New(ctx.Context(), entities.Session{
		UserId:       userId,
		RefreshToken: security.generator.Refresh(),
		UserAgent:    string(ctx.Context().UserAgent()),
		Ip:           ctx.IP(),
		ExpiresIn:    time.Now().AddDate(0, 0, 30),
	})
	if err != nil {
		return "", err
	}
	return newSession.RefreshToken, nil
}

// Check refresh token in the black list. If contains then return true
func (security *Security) IsSessionBanned(ctx context.Context, refreshToken string) bool {
	err := security.sessions.CheckIfExist(ctx, refreshToken)
	if err != nil {
		return false
	}
	return true
}

// Adding tokens to black list, so that the user can not access service with
// passed refresh token
func (security *Security) BanSessions(ctx context.Context, tokens ...entities.Session) error {
	for _, t := range tokens {
		if err := security.sessions.Add(
			ctx,
			t.RefreshToken,
			t.ID,
			time.Duration(t.ExpiresIn.Unix()*1000),
		); err != nil {
			return err
		}
	}
	return nil
}

// SignAccessToken create new JWTResponse for user
//
// If withSession flag passed additionally create refresh session
func (security *Security) SignAccessToken(ctx *fiber.Ctx, refreshToken string) (JWTResponse, error) {
	var (
		jwt JWTResponse
		err error
	)
	// Get user from context
	claims, ok := ctx.Locals("request_user").(*UserClaims)
	if !ok {
		return JWTResponse{}, ErrEmptyContext
	}
	jwt.TokenType = "Bearer"
	jwt.ExpiresIn = int(security.config.JwtExpires)
	// If withSession additionally create an refresh session
	if refreshToken != "" {
		jwt.RefreshToken = refreshToken
	}
	// Creating access token
	jwt.AccessToken, err = security.generator.JwtWithRefresh(claims.ID, jwt.RefreshToken)
	if err != nil {
		return JWTResponse{}, err
	}

	return jwt, nil
}

func (security *Security) ValidateJwt(token string) (*Claims, error) {
	return security.generator.ValidateJwt(token)
}

func NewSecurity(config config.Authorization, sessions sessions.Sessions) *Security {
	return &Security{
		config:    config,
		generator: newSecurityGen(config),
		sessions:  sessions,
	}
}
