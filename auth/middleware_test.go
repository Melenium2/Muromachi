package auth_test

import (
	"Muromachi/auth"
	"Muromachi/config"
	"Muromachi/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestApplyAuthMiddleware_Mock(t *testing.T) {
	app := fiber.New()

	cfg := config.Authorization{
		JwtSalt:    "hiprivetkonichiua",
		JwtExpires: time.Hour * 1,
		JwtIss:     "auth.apptwice.com",
		JwtAud:     "",
	}
	defender := auth.NewSecurity(cfg, mockSession{})
	app.Use("/", auth.ApplyAuthMiddleware(defender))
	app.Get("/test", func(ctx *fiber.Ctx) error {
		s, _ := mockSession{}.Get(ctx.Context(), "123")
		return ctx.JSON(s)
	})

	var tt = []struct {
		name               string
		withJwt            bool
		withNotValidJwt    bool
		withExpiredJwt     bool
		withCookieJwt      bool
		expectedStatusCode int
	}{
		{
			name:               "401 error, request without jwt",
			withJwt:            false,
			withCookieJwt:      false,
			expectedStatusCode: 401,
		},
		{
			name:               "200 if jwt found in header values",
			withJwt:            true,
			withCookieJwt:      false,
			expectedStatusCode: 200,
		},
		{
			name:               "200 if jwt found in cookies values",
			withJwt:            false,
			withCookieJwt:      true,
			expectedStatusCode: 200,
		},
		{
			name:               "401 if jwt token not valid",
			withJwt:            true,
			withNotValidJwt:    true,
			withCookieJwt:      true,
			expectedStatusCode: 401,
		},
		{
			name:               "401 if expired jwt",
			withJwt:            true,
			withExpiredJwt:     true,
			expectedStatusCode: 401,
		},
	}

	// TODO Добавить кейс с isBaned
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			localTime := time.Now().Add(cfg.JwtExpires)
			if test.withExpiredJwt {
				localTime = localTime.Add(time.Hour * -24)
			}
			rtoken, _ := utils.UUID()
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
				StandardClaims: &jwt.StandardClaims{
					Audience:  cfg.JwtAud,
					ExpiresAt: localTime.Unix(),
					IssuedAt:  time.Now().Unix(),
					Issuer:    cfg.JwtIss,
					Id:        rtoken,
				},
				UserClaims: &auth.UserClaims{
					ID:   123,
					Role: "user",
				},
			})
			token, err := jwtToken.SignedString([]byte(cfg.JwtSalt))
			assert.NoError(t, err)

			if test.withNotValidJwt {
				token += "123"
			}
			req, _ := http.NewRequest("GET", "/test", nil)
			if test.withJwt {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}
			if test.withCookieJwt {
				req.AddCookie(&http.Cookie{
					Name:       auth.SecurityCookieName,
					Value:      token,
					Path:       "/",
					RawExpires: "",
					MaxAge:     int(time.Now().AddDate(0, 0, 1).Unix()),
					HttpOnly:   true,
				})
			}

			resp, err := app.Test(req, 1000*60)
			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)

			b, _ := ioutil.ReadAll(resp.Body)
			assert.NotNil(t, b)
			t.Log("msg", " ", string(b))
		})
	}
}
