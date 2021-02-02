package auth_test

import (
	"Muromachi/auth"
	"Muromachi/config"
	"Muromachi/store/entities"
	"Muromachi/store/testhelpers"
	"Muromachi/store/users/sessions"
	"Muromachi/store/users/sessions/blacklist"
	"Muromachi/store/users/sessions/tokens"
	"Muromachi/store/users/userstore"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
	"time"
)

func TestSecurity_StartSession_Mock(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "hiprivetsalt",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	var tt = []struct {
		name          string
		session       sessions.Session
		refreshToken  string
		passUserCtx   bool
		expectedError bool
	}{
		{
			name:          "create new instance of refresh token",
			session:       mockSession{},
			refreshToken:  "",
			passUserCtx:   true,
			expectedError: false,
		},
		{
			name:          "can not find session with given id",
			session:       mockSessionRemoveNoRows{},
			refreshToken:  "123",
			passUserCtx:   false,
			expectedError: true,
		},
		{
			name:          "expired refresh session",
			session:       mockSessionRemoveExpiredSession{},
			refreshToken:  "123",
			passUserCtx:   false,
			expectedError: true,
		},
		{
			name:          "session in black list",
			session:       mockSessionBannedToken{},
			refreshToken:  "123",
			passUserCtx:   false,
			expectedError: true,
		},
		{
			name:          "user has more then 5 opened sessions",
			session:       mockSessionMoreThen5{},
			refreshToken:  "",
			passUserCtx:   true,
			expectedError: false,
		},
		{
			name:          "don't not how but ctx value is empty",
			session:       mockSession{},
			refreshToken:  "",
			passUserCtx:   false,
			expectedError: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			security := auth.NewSecurity(cfg, test.session)
			app := fiber.New()
			fastCtx := &fasthttp.RequestCtx{}
			ctx := app.AcquireCtx(fastCtx)

			if test.passUserCtx {
				ctx.Locals("request_user", &auth.UserClaims{
					ID:   123,
					Role: "user",
				})
			}
			token, err := security.StartSession(ctx, test.refreshToken)
			assert.Equal(t, err != nil, test.expectedError)
			if !test.expectedError {
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestSecurity_StartSession_ShouldCreateNewRefreshSession(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "hiprivetsalt",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}
	dbcfg := config.New("../config/dev.yml").Database
	dbcfg.Schema = "../config/schema.sql"
	conn, cleaner := testhelpers.RealDb(dbcfg)
	defer cleaner("refresh_sessions", "users")

	sess := tokens.New(conn)
	security := auth.NewSecurity(cfg, sessions.New(sess, nil))

	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	user, _ := userstore.NewUserRepo(conn).Create(context.Background(), u)

	// FIX
	// Bad solution. But i have troubles with fasthttp context.
	// Don't know why, but pgx freezes with empty context.
	// Or need emulate working fasthttp ctx.
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		ctx.Locals("request_user", &auth.UserClaims{
			ID:   int64(user.ID),
			Role: "user",
		})

		token, err := security.StartSession(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		return nil
	})

	req, _ := http.NewRequest("GET", "/", nil)
	_, _ = app.Test(req)
}

func TestSecurity_StartSession_ShouldReturnErrorIfRefreshSessionNotFound(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "hiprivetsalt",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}
	dbcfg := config.New("../config/dev.yml").Database
	dbcfg.Schema = "../config/schema.sql"
	conn, cleaner := testhelpers.RealDb(dbcfg)
	defer cleaner("refresh_sessions")

	sess := tokens.New(conn)
	security := auth.NewSecurity(cfg, sessions.New(sess, nil))

	// FIX
	// Bad solution. But i have troubles with fasthttp context.
	// Don't know why, but pgx freezes with empty context.
	// Or need emulate working fasthttp ctx.
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {

		token, err := security.StartSession(ctx, "123")
		assert.Error(t, err)
		assert.Empty(t, token)

		return nil
	})

	req, _ := http.NewRequest("GET", "/", nil)
	_, _ = app.Test(req)
}

func TestSecurity_StartSession_ShouldReturnErrorIfRefreshSessionIsExpired(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "hiprivetsalt",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}
	dbcfg := config.New("../config/dev.yml").Database
	dbcfg.Schema = "../config/schema.sql"
	conn, cleaner := testhelpers.RealDb(dbcfg)
	defer cleaner("refresh_sessions", "users")

	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	user, _ := userstore.NewUserRepo(conn).Create(context.Background(), u)

	sess := tokens.New(conn)
	_, _ = sess.New(context.Background(), entities.Session{
		UserId:       user.ID,
		RefreshToken: "123",
		UserAgent:    "",
		Ip:           "",
		ExpiresIn:    time.Now().Add(time.Hour * -24),
	})

	security := auth.NewSecurity(cfg, sessions.New(sess, nil))

	// FIX
	// Bad solution. But i have troubles with fasthttp context.
	// Don't know why, but pgx freezes with empty context.
	// Or need to simulate working fasthttp ctx.
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {

		token, err := security.StartSession(ctx, "123")
		assert.Error(t, err)
		assert.Empty(t, token)

		return nil
	})

	req, _ := http.NewRequest("GET", "/", nil)
	_, _ = app.Test(req)
}

func TestSecurity_StartSession_ShouldReturnErrorIfSessionInBlackList(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "hiprivetsalt",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}
	dbcfg := config.New("../config/dev.yml")
	dbcfg.Database.Schema = "../config/schema.sql"
	conn, cleaner := testhelpers.RealDb(dbcfg.Database)
	defer cleaner("refresh_sessions", "users")

	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	user, _ := userstore.NewUserRepo(conn).Create(context.Background(), u)

	sess := tokens.New(conn)
	_, _ = sess.New(context.Background(), entities.Session{
		UserId:       user.ID,
		RefreshToken: "123",
		UserAgent:    "",
		Ip:           "",
		ExpiresIn:    time.Now().Add(time.Hour * 24),
	})

	redisConn, redisCleaner := testhelpers.RedisDb(dbcfg.Database.Redis)
	defer redisCleaner()

	balcklist := blacklist.New(redisConn)
	assert.NoError(t, balcklist.Add(context.Background(), "123", 1, time.Hour))

	security := auth.NewSecurity(cfg, sessions.New(sess, balcklist))

	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		token, err := security.StartSession(ctx, "123")
		assert.Error(t, err)
		t.Log(err)
		assert.Empty(t, token)

		return nil
	})

	req, _ := http.NewRequest("GET", "/", nil)
	_, _ = app.Test(req)
}

func TestSecurity_StartSession_ShouldReturnErrorIfRefreshTokenNotProvidedAndContextValuesEmpty(t *testing.T) {
	security := auth.NewSecurity(config.Authorization{}, sessions.New(nil, nil))
	app := fiber.New()
	fastCtx := &fasthttp.RequestCtx{}
	ctx := app.AcquireCtx(fastCtx)
	token, err := security.StartSession(ctx)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestSecurity_StartSession_ShouldRemoveAllUserSessionsIfLenIsMoreThen5AndReturnNewToken(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "hiprivetsalt",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}
	dbcfg := config.New("../config/dev.yml").Database
	dbcfg.Schema = "../config/schema.sql"
	conn, cleaner := testhelpers.RealDb(dbcfg)
	defer cleaner("refresh_sessions", "users")

	u := entities.User{
		Company: "123",
	}
	_ = u.GenerateSecrets()
	user, _ := userstore.NewUserRepo(conn).Create(context.Background(), u)

	sess := tokens.New(conn)
	for i := 0; i < 7; i++ {
		_, _ = sess.New(context.Background(), entities.Session{
			UserId:       user.ID,
			RefreshToken: fmt.Sprintf("123%d", i+1),
			UserAgent:    "",
			Ip:           "",
			ExpiresIn:    time.Now().Add(time.Hour * 24),
		})
	}

	security := auth.NewSecurity(cfg, sessions.New(sess, nil))

	// FIX
	// Bad solution. But i have troubles with fasthttp context.
	// Don't know why, but pgx freezes with empty context.
	// Or need to simulate working fasthttp ctx.
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		ctx.Locals("request_user", &auth.UserClaims{
			ID:   int64(user.ID),
			Role: "user",
		})

		token, err := security.StartSession(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		return nil
	})

	req, _ := http.NewRequest("GET", "/", nil)
	_, _ = app.Test(req)

	ses, err := sess.UserSessions(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ses))
}

func TestSecurity_SignAccessToken(t *testing.T) {
	cfg := config.Authorization{
		JwtSalt:    "uuuuuuuuuuuuuuuuuuthen",
		JwtExpires: time.Hour * 24,
		JwtIss:     "apptwice.com",
	}

	var tt = []struct {
		name          string
		withUserCtx   bool
		expectedError bool
	}{
		{
			name:          "error if user not in context",
			withUserCtx:   false,
			expectedError: true,
		},
		{
			name:          "create new access token",
			withUserCtx:   true,
			expectedError: false,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			security := auth.NewSecurity(cfg, nil)

			app := fiber.New()
			fastCtx := &fasthttp.RequestCtx{}
			ctx := app.AcquireCtx(fastCtx)

			if test.withUserCtx {
				ctx.Locals("request_user", &auth.UserClaims{
					ID:   123,
					Role: "user",
				})
			}

			token, err := security.SignAccessToken(ctx, "")
			assert.Equal(t, err != nil, test.expectedError)
			if test.expectedError {
				assert.Empty(t, token)
			} else {
				assert.NotEmpty(t, token)
			}
		})
	}
}
