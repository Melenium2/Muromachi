package server_test

import (
	"Muromachi/auth"
	"Muromachi/config"
	"Muromachi/server"
	"Muromachi/server/requests"
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
	"fmt"
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
	// Mock config for authorization process
	cfg := config.Authorization{
		JwtSalt:    "nunetprivet",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	// Pepare handler
	sec := auth.NewSecurity(cfg, mockSession{})
	col := users.NewAuthTables(mockSession{}, userstore.NewUserRepo(mockConn{}))

	handler := server.Authorize(sec, col)

	// Prepare fiber app
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
			// Encode data with json encoding or x-www-form-urlencoded
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
	// Mock config for authorization process
	cfg := config.Authorization{
		JwtSalt:    "nunetprivet",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	// Prepare handler
	sec := auth.NewSecurity(cfg, mockSession{})
	col := users.NewAuthTables(mockSession{}, userstore.NewUserRepo(mockConn{}))

	handler := server.Authorize(sec, col)

	app := fiber.New()
	app.Post("/auth", handler)

	// Generate wrong data type
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
	// Load config from file
	cfg := config.New("../config/dev.yml")
	// Change path to database schema
	cfg.Database.Schema = "../config/schema.sql"
	// Conn to real db
	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users", "refresh_sessions")

	// Add auth config
	cfg.Auth = config.Authorization{
		JwtSalt:    "nunetprivet",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	// Prepare db
	userRepo := userstore.NewUserRepo(conn)
	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	clientId, clientSecret := u.ClientId, u.ClientSecret
	user, err := userRepo.Create(context.Background(), u)
	assert.NoError(t, err)

	// Conn to redis
	redisConn, redisCleaner := testhelpers.RedisDb(cfg.Database.Redis)
	defer redisCleaner()
	blacklist := blacklist.New(redisConn)

	// Prepare handler
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
			// For each test generate 2 sessions
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

			// Add to black list if true
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

func TestBan_Mock(t *testing.T) {
	tables := users.NewAuthTables(mockSession{}, userstore.NewUserRepo(mockConn{}))

	handler := server.Ban(tables)

	app := fiber.New()
	app.Post("/ban", handler)

	var tt = []struct {
		name              string
		tokenList         requests.TokenList
		withWrongDataType interface{}
		expectedCode      int
		expectedBans      int
	}{
		{
			name: "should ban both user and refresh tokens",
			tokenList: requests.TokenList{
				UserId: 13,
				Tokens: []string{"123", "321", "132"},
				Ttl:    int(time.Hour * 24),
			},
			expectedCode: 200,
			expectedBans: 5,
		},
		{
			name:              "should return error if wrong data type",
			withWrongDataType: []int{1, 2, 3},
			expectedCode:      400,
		},
		{
			name: "should ban only user session",
			tokenList: requests.TokenList{
				UserId: 15,
				Tokens: nil,
				Ttl:    int(time.Hour * 24),
			},
			expectedCode: 200,
			expectedBans: 2,
		},
		{
			name: "should ban only refresh tokens",
			tokenList: requests.TokenList{
				Tokens: []string{"123", "321", "132"},
				Ttl:    int(time.Hour * 24),
			},
			expectedCode: 200,
			expectedBans: 3,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			var (
				b  io.Reader
				by []byte
			)
			// Prepare data. If withWrongDataType is not nil encode wrong data type
			if test.withWrongDataType != nil {
				by, _ = json.Marshal(test.withWrongDataType)
			} else {
				by, _ = json.Marshal(test.tokenList)
			}
			b = bytes.NewBuffer(by)

			req, _ := http.NewRequest("POST", "/ban", b)
			req.Header.Set("content-type", "application/json")
			resp, err := app.Test(req, 1000*60)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			// Check count
			var info requests.BanInfo
			by, _ = ioutil.ReadAll(resp.Body)
			_ = json.Unmarshal(by, &info)
			assert.Equal(t, test.expectedBans, info.Count)
		})
	}
}

func TestBan(t *testing.T) {
	// Load config
	cfg := config.New("../config/dev.yml")
	// Change path to database schema
	cfg.Database.Schema = "../config/schema.sql"
	// Conn to real db
	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users", "refresh_sessions")

	// Prepare db
	userRepo := userstore.NewUserRepo(conn)
	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	user, err := userRepo.Create(context.Background(), u)
	assert.NoError(t, err)

	// Conn to redis
	redisConn, redisCleaner := testhelpers.RedisDb(cfg.Database.Redis)
	defer redisCleaner()
	blacklist := blacklist.New(redisConn)

	// Prepare handler
	sessionRepo := sessions.New(tokens.New(conn), blacklist)
	col := users.NewAuthTables(sessionRepo, userRepo)

	handler := server.Ban(col)

	var tt = []struct {
		name              string
		tokenList         requests.TokenList
		withWrongDataType interface{}
		expectedCode      int
		expectedBans      int
		// Like before each and after each
		precomputed func() func()
	}{
		{
			name: "should ban both user and refresh tokens",
			tokenList: requests.TokenList{
				UserId: user.ID,
				Tokens: []string{"123"},
				Ttl:    int(time.Hour * 24),
			},
			expectedCode: 200,
			expectedBans: 6,
			precomputed: func() func() {
				ses := entities.Session{UserId: user.ID, RefreshToken: "123"}
				for i := 0; i < 5; i++ {
					_, err = sessionRepo.New(context.Background(), ses)
					assert.NoError(t, err)
					ses.RefreshToken += fmt.Sprintf("%d", i+1)
				}

				return func() {
					cleaner("refresh_sessions")
					redisCleaner()
				}
			},
		},
		{
			name: "should return pgx error no rows in result set",
			tokenList: requests.TokenList{
				Tokens: []string{"123"},
			},
			expectedCode: 400,
		},
		{
			name:              "should return error if wrong data type",
			withWrongDataType: []int{1, 2, 3},
			expectedCode:      400,
		},
		{
			name: "should ban only user session",
			tokenList: requests.TokenList{
				UserId: user.ID,
				Ttl:    int(time.Hour * 24),
			},
			expectedCode: 200,
			expectedBans: 3,
			precomputed: func() func() {
				ses := entities.Session{UserId: user.ID, RefreshToken: "123"}
				for i := 0; i < 3; i++ {
					_, err = sessionRepo.New(context.Background(), ses)
					assert.NoError(t, err)
					ses.RefreshToken += fmt.Sprintf("%d", i+1)
				}

				return func() {
					cleaner("refresh_sessions")
					redisCleaner()
				}
			},
		},
		{
			name: "should ban only refresh tokens",
			tokenList: requests.TokenList{
				Tokens: []string{"123", "321", "132"},
				Ttl:    int(time.Hour * 24),
			},
			expectedCode: 200,
			expectedBans: 3,
			precomputed: func() func() {
				ses := entities.Session{UserId: user.ID, RefreshToken: "123"}
				_, err = sessionRepo.New(context.Background(), ses)
				ses.RefreshToken = "321"
				_, err = sessionRepo.New(context.Background(), ses)
				ses.RefreshToken = "132"
				_, err = sessionRepo.New(context.Background(), ses)

				return func() {
					cleaner("refresh_sessions")
					redisCleaner()
				}
			},
		},
	}

	app := fiber.New()
	app.Post("/ban", handler)

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			var (
				b     io.Reader
				by    []byte
				clean func()
			)
			// If not nil before each test run precomputed()
			if test.precomputed != nil {
				clean = test.precomputed()
			}
			if test.withWrongDataType != nil {
				by, _ = json.Marshal(test.withWrongDataType)
			} else {
				by, _ = json.Marshal(test.tokenList)
			}
			b = bytes.NewReader(by)

			req, _ := http.NewRequest("POST", "/ban", b)
			req.Header.Set("content-type", "application/json")
			resp, err := app.Test(req, 1000*60)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			var info requests.BanInfo
			by, _ = ioutil.ReadAll(resp.Body)
			t.Log(string(by))
			_ = json.Unmarshal(by, &info)
			assert.Equal(t, test.expectedBans, info.Count)

			// if clean func not nil then run func
			if clean != nil {
				clean()
			}
		})
	}
}

func TestUnban_Mock(t *testing.T) {
	tables := users.NewAuthTables(mockSession{}, userstore.NewUserRepo(mockConn{}))

	handler := server.Ban(tables)

	app := fiber.New()
	app.Post("/ban", handler)

	var tt = []struct {
		name              string
		tokenList         requests.TokenList
		withWrongDataType interface{}
		expectedCode      int
		expectedDeletions int
	}{
		{
			name: "should unban user and refresh tokens",
			tokenList: requests.TokenList{
				UserId: 123,
				Tokens: []string{"321", "123", "444"},
			},
			expectedCode:      200,
			expectedDeletions: 5,
		},
		{
			name:              "should return error if incorrect datatype",
			withWrongDataType: []int{123, 123, 123},
			expectedCode:      400,
		},
		{
			name: "should unban all user sessions",
			tokenList: requests.TokenList{
				UserId: 123,
			},
			expectedCode:      200,
			expectedDeletions: 2,
		},
		{
			name: "should unban given refresh tokens",
			tokenList: requests.TokenList{
				Tokens: []string{"123"},
			},
			expectedCode:      200,
			expectedDeletions: 1,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			var (
				b  io.Reader
				by []byte
			)
			if test.withWrongDataType != nil {
				by, _ = json.Marshal(test.withWrongDataType)
			} else {
				by, _ = json.Marshal(test.tokenList)
			}
			b = bytes.NewBuffer(by)

			req, _ := http.NewRequest("POST", "/ban", b)
			req.Header.Set("content-type", "application/json")
			resp, err := app.Test(req, 1000*60)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			var info requests.BanInfo
			by, _ = ioutil.ReadAll(resp.Body)
			_ = json.Unmarshal(by, &info)
			assert.Equal(t, test.expectedDeletions, info.Count)
		})
	}
}

func TestUnban(t *testing.T) {
	// Load config
	cfg := config.New("../config/dev.yml")
	// Change path to database schema
	cfg.Database.Schema = "../config/schema.sql"
	// Conn to real db
	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users", "refresh_sessions")

	ctx := context.Background()

	// Prepare db
	userRepo := userstore.NewUserRepo(conn)
	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	user, err := userRepo.Create(context.Background(), u)
	assert.NoError(t, err)

	// Conn to redis
	redisConn, redisCleaner := testhelpers.RedisDb(cfg.Database.Redis)
	defer redisCleaner()
	blacklist := blacklist.New(redisConn)

	// Prepare handler
	sessionRepo := sessions.New(tokens.New(conn), blacklist)
	col := users.NewAuthTables(sessionRepo, userRepo)

	handler := server.Unban(col)

	var tt = []struct {
		name              string
		tokenList         requests.TokenList
		withWrongDataType interface{}
		expectedCode      int
		expectedDeletions int
		// Like before each and after each
		precomputed func() func()
	}{
		{
			name: "should unban user and refresh tokens",
			tokenList: requests.TokenList{
				UserId: user.ID,
				Tokens: []string{"000"},
			},
			expectedCode:      200,
			expectedDeletions: 4,
			precomputed: func() func() {
				l := []string{"321", "123", "444"}
				for _, token := range l {
					s, _ := sessionRepo.New(ctx, entities.Session{UserId: user.ID, RefreshToken: token})
					assert.NoError(t, sessionRepo.Add(context.Background(), s.RefreshToken, s.ID, time.Hour))
				}
				_ = sessionRepo.Add(context.Background(), "000", 2, time.Hour)

				return func() {
					cleaner("refresh_sessions")
					redisCleaner()
				}
			},
		},
		{
			name:              "should return error if incorrect datatype",
			withWrongDataType: []int{123, 123, 123},
			expectedCode:      400,
		},
		{
			name: "should unban all user sessions",
			tokenList: requests.TokenList{
				UserId: user.ID,
			},
			expectedCode:      200,
			expectedDeletions: 3,
			precomputed: func() func() {
				l := []string{"321", "123", "444"}
				for _, token := range l {
					s, _ := sessionRepo.New(ctx, entities.Session{UserId: user.ID, RefreshToken: token})
					assert.NoError(t, sessionRepo.Add(context.Background(), s.RefreshToken, s.ID, time.Hour))
				}

				return func() {
					cleaner("refresh_sessions")
					redisCleaner()
				}
			},
		},
		{
			name: "should unban given refresh tokens",
			tokenList: requests.TokenList{
				Tokens: []string{"123"},
			},
			expectedCode:      200,
			expectedDeletions: 1,
			precomputed: func() func() {
				l := []string{"123"}
				for _, token := range l {
					s, _ := sessionRepo.New(ctx, entities.Session{UserId: user.ID, RefreshToken: token})
					assert.NoError(t, sessionRepo.Add(context.Background(), s.RefreshToken, s.ID, time.Hour))
				}

				return func() {
					cleaner("refresh_sessions")
					redisCleaner()
				}
			},
		},
	}

	app := fiber.New()
	app.Post("/ban", handler)

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			var (
				b     io.Reader
				by    []byte
				clean func()
			)
			// if precomputed not nil then run before test
			if test.precomputed != nil {
				clean = test.precomputed()
			}
			if test.withWrongDataType != nil {
				by, _ = json.Marshal(test.withWrongDataType)
			} else {
				by, _ = json.Marshal(test.tokenList)
			}
			b = bytes.NewReader(by)

			req, _ := http.NewRequest("POST", "/ban", b)
			req.Header.Set("content-type", "application/json")
			resp, err := app.Test(req, 1000*60)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			var info requests.BanInfo
			by, _ = ioutil.ReadAll(resp.Body)
			_ = json.Unmarshal(by, &info)
			assert.Equal(t, test.expectedDeletions, info.Count)

			// if clean not nil then run after test
			if clean != nil {
				clean()
			}
		})
	}
}
