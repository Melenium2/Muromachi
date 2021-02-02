package server_test

import (
	"Muromachi/auth"
	"Muromachi/config"
	"Muromachi/server"
	"Muromachi/store/entities"
	"Muromachi/store/testhelpers"
	"Muromachi/store/users"
	"Muromachi/store/users/sessions"
	"Muromachi/store/users/sessions/blacklist"
	"Muromachi/store/users/sessions/tokens"
	"Muromachi/store/users/userstore"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestAuthorize_Mock(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "nunetprivet",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	sec := auth.NewSecurity(cfg, mockSession{})
	col := users.NewAuthTables(mockSession{}, userstore.NewUserRepo(mockConn{}))

	handler := server.Authorize(sec, col)

	app := fiber.New()
	app.Post("/auth", handler)

	var tt = []struct {
		name          string
		withJson      bool
		withUrlValues bool
		request       auth.JWTRequest
		expectedCode  int
	}{
		{
			name:          "request with form, and right data, should return new access token",
			withJson:      false,
			withUrlValues: true,
			request: auth.JWTRequest{
				AccessType:   "simple",
				ClientId:     "123",
				ClientSecret: "123",
			},
			expectedCode: 200,
		},
		{
			name:          "request with json, and right data, should return new access token",
			withJson:      true,
			withUrlValues: false,
			request: auth.JWTRequest{
				AccessType:   "simple",
				ClientId:     "123",
				ClientSecret: "123",
			},
			expectedCode: 200,
		},
		{
			name:          "request with wrong user client id, should return 404 error",
			withJson:      true,
			withUrlValues: false,
			request: auth.JWTRequest{
				AccessType:   "simple",
				ClientId:     "1234",
				ClientSecret: "123",
			},
			expectedCode: 404,
		},
		{
			name:          "request with wrong user client secret, should return 401 error",
			withJson:      false,
			withUrlValues: true,
			request: auth.JWTRequest{
				AccessType:   "simple",
				ClientId:     "123",
				ClientSecret: "1234",
			},
			expectedCode: 401,
		},
		{
			name:          "request with new session, should return 200 and new access token",
			withJson:      false,
			withUrlValues: true,
			request: auth.JWTRequest{
				AccessType:   "session",
				ClientId:     "123",
				ClientSecret: "123",
			},
			expectedCode: 200,
		},
		{
			name:          "request for update refresh token, should return 200 and new access token and refresh token",
			withJson:      false,
			withUrlValues: true,
			request: auth.JWTRequest{
				AccessType:   "refresh_token",
				ClientId:     "123",
				ClientSecret: "123",
			},
			expectedCode: 200,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			var (
				b           io.Reader
				contentType string
			)
			if test.withJson {
				by, _ := json.Marshal(test.request)
				b = bytes.NewReader(by)
				contentType = "application/json"
			}
			if test.withUrlValues {
				data := url.Values{}
				data.Set("access_type", test.request.AccessType)
				data.Set("client_id", test.request.ClientId)
				data.Set("client_secret", test.request.ClientSecret)
				data.Set("refresh_token", test.request.RefreshToken)
				b = strings.NewReader(data.Encode())
				contentType = "application/x-www-form-urlencoded"
			}
			req, _ := http.NewRequest("POST", "/auth", b)
			req.Header.Set("Content-Type", contentType)

			resp, err := app.Test(req, 1000*60)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			by, _ := ioutil.ReadAll(resp.Body)
			t.Log(string(by))
		})
	}
}

func TestAuthorize_WithWrongRequestDataType(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "nunetprivet",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	sec := auth.NewSecurity(cfg, mockSession{})
	col := users.NewAuthTables(mockSession{}, userstore.NewUserRepo(mockConn{}))

	handler := server.Authorize(sec, col)

	app := fiber.New()
	app.Post("/auth", handler)

	by, _ := json.Marshal(map[string]interface{}{
		"hihi": "123",
	})
	b := bytes.NewReader(by)
	req, _ := http.NewRequest("POST", "/auth", b)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 1000*60)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode)

	by, _ = ioutil.ReadAll(resp.Body)
	t.Log(string(by))
}

func TestAuthorize(t *testing.T) {
	cfg := config.New("../config/dev.yml")
	cfg.Database.Schema = "../config/schema.sql"
	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users", "refresh_sessions")

	cfg.Auth = config.Authorization{
		JwtSalt:    "nunetprivet",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	userRepo := userstore.NewUserRepo(conn)
	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	clientId, clientSecret := u.ClientId, u.ClientSecret
	user, err := userRepo.Create(context.Background(), u)
	assert.NoError(t, err)

	redisConn, redisCleaner := testhelpers.RedisDb(cfg.Database.Redis)
	defer redisCleaner()
	blacklist := blacklist.New(redisConn)

	sessionRepo := sessions.New(tokens.New(conn), blacklist)

	sec := auth.NewSecurity(cfg.Auth, sessionRepo)
	col := users.NewAuthTables(sessionRepo, userRepo)

	handler := server.Authorize(sec, col)

	app := fiber.New()
	app.Post("/auth", handler)

	var tt = []struct {
		name         string
		inBlackList  bool
		request      auth.JWTRequest
		expectedCode int
	}{
		{
			name: "should return new access token",
			request: auth.JWTRequest{
				AccessType:   "simple",
				ClientId:     clientId,
				ClientSecret: clientSecret,
			},
			expectedCode: 200,
		},
		{
			name: "should return new access token with access type session",
			request: auth.JWTRequest{
				AccessType:   "session",
				ClientId:     clientId,
				ClientSecret: clientSecret,
			},
			expectedCode: 200,
		},
		{
			name: "should update refresh token",
			request: auth.JWTRequest{
				AccessType:   "refresh_token",
				ClientId:     clientId,
				ClientSecret: clientSecret,
				RefreshToken: "123",
			},
			expectedCode: 200,
		},
		{
			name: "should return error if refresh token is expired",
			request: auth.JWTRequest{
				AccessType:   "refresh_token",
				ClientId:     clientId,
				ClientSecret: clientSecret,
				RefreshToken: "321",
			},
			expectedCode: 400,
		},
		{
			name: "should return error if client id not found",
			request: auth.JWTRequest{
				AccessType:   "simple",
				ClientId:     clientId + "123123123",
				ClientSecret: clientSecret,
			},
			expectedCode: 404,
		},
		{
			name: "should return error if client secret not valid",
			request: auth.JWTRequest{
				AccessType:   "simple",
				ClientId:     clientId,
				ClientSecret: clientSecret + "1",
			},
			expectedCode: 401,
		},
		{
			name: "should return error if access type is wrong",
			request: auth.JWTRequest{
				AccessType:   "wrong",
				ClientId:     clientId,
				ClientSecret: clientSecret,
			},
			expectedCode: 404,
		},
		{
			name:        "should return error if refresh token in blacklist",
			inBlackList: true,
			request: auth.JWTRequest{
				AccessType:   "refresh_token",
				RefreshToken: "123",
				ClientId:     clientId,
				ClientSecret: clientSecret,
			},
			expectedCode: 400,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			_, _ = sessionRepo.New(context.Background(), entities.Session{
				UserId:       user.ID,
				RefreshToken: "123",
				UserAgent:    "",
				Ip:           "",
				ExpiresIn:    time.Now().AddDate(0, 0, 1),
			})
			_, _ = sessionRepo.New(context.Background(), entities.Session{
				UserId:       user.ID,
				RefreshToken: "321",
				UserAgent:    "",
				Ip:           "",
				ExpiresIn:    time.Now().AddDate(0, 0, -1),
			})
			defer cleaner("refresh_sessions")

			by, _ := json.Marshal(test.request)
			b := bytes.NewReader(by)

			req, _ := http.NewRequest("POST", "/auth", b)
			req.Header.Set("Content-Type", "application/json")

			if test.inBlackList {
				_ = blacklist.Add(context.Background(), test.request.RefreshToken, "123", time.Hour)
			}

			resp, err := app.Test(req, 1000*60)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			by, _ = ioutil.ReadAll(resp.Body)
			t.Log(string(by))
		})
	}

}
