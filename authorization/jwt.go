package authorization

import (
	"Muromachi/config"
	"Muromachi/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var (
	ErrUnexpectedJwtError  = errors.New("unexpected jwt error")
	ErrInvalidToken        = errors.New("invalid jwt token")
	ErrExpiredAccessToken  = errors.New("expired access token")
	ErrExpiredRefreshToken = errors.New("expired refresh token")
	ErrEmptyContext        = errors.New("empty context")
)

type UserClaims struct {
	ID   int64
	Role string
}

type Claims struct {
	*jwt.StandardClaims
	*UserClaims
}

type securityGenerator struct {
	config config.Authorization
}

func (gen *securityGenerator) Refresh() string {
	uuid, err := utils.UUID()
	if err != nil {
		return ""
	}
	return utils.Hash(uuid, time.Now().Unix())
}

func (gen *securityGenerator) Jwt(userId int64) (string, error) {
	return gen.JwtWithRefresh(userId, gen.Refresh())
}

func (gen *securityGenerator) JwtWithRefresh(userId int64, refreshToken string) (string, error) {
	if gen.config.JwtSalt == "" {
		log.Fatal("jwt salt env not provided")
	}
	t := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: &jwt.StandardClaims{
			Audience:  gen.config.JwtAud,
			ExpiresAt: t.Add(gen.config.JwtExpires).Unix(),
			IssuedAt:  t.Unix(),
			Issuer:    gen.config.JwtIss,
			Id:        refreshToken,
		},
		UserClaims: &UserClaims{
			ID:   userId,
			Role: "userrepo",
		},
	})

	return token.SignedString(gen.config.JwtSalt)
}

func (gen *securityGenerator) ValidateJwt(token string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(gen.config.JwtSalt), nil
	})
	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil, ErrUnexpectedJwtError
	}

	return claims, nil
}

func newSecurityGen(config config.Authorization) *securityGenerator {
	return &securityGenerator{
		config: config,
	}
}
