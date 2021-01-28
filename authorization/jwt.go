package authorization

import (
	"Muromachi/config"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
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

func (gen *securityGenerator) UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func (gen *securityGenerator) Hash(str string, random interface{}) string {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		panic(fmt.Sprintf("terrible hash, %v", err))
	}
	secret := fmt.Sprintf("%s_%s_%v", str, salt, random)
	hash := sha256.Sum256([]byte(secret))

	return fmt.Sprintf("%x", hash)
}

func (gen *securityGenerator) Refresh() string {
	uuid, err := gen.UUID()
	if err != nil {
		return ""
	}
	return gen.Hash(uuid, time.Now().Unix())
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
			Role: "user",
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
