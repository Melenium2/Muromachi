package sessionrepo_test

import (
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"Muromachi/store/sessionrepo"
	"Muromachi/store/testhelpers"
	user2 "Muromachi/store/userrepo"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRefreshRepo_New(t *testing.T) {
	conn, cleaner := testhelpers.RealDb()
	defer cleaner("users")

	repo := user2.NewUserRepo(conn)
	user := entities.User{Company: "123"}
	_ = user.GenerateSecrets()

	u, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	var tt = []struct {
		name          string
		conn          connector.Conn
		cleaner       func(s ...string)
		expectedError bool
		doError       bool
	}{
		{
			name:          "mock | should create new userrepo",
			conn:          mockRefreshNewFuncConnSuccess{},
			cleaner:       nil,
			expectedError: false,
		},
		{
			name:          "should add new userrepo to db",
			conn:          conn,
			cleaner:       cleaner,
			expectedError: false,
		},
		{
			name:          "mock | should return error if can not add new userrepo",
			conn:          mockRefreshNewFuncConnError{},
			cleaner:       nil,
			expectedError: true,
		},
		{
			name:          "should return err because connection was closed =)",
			conn:          conn,
			cleaner:       nil,
			expectedError: true,
			doError:       true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if test.cleaner != nil {
				defer cleaner("refresh_sessions")
			}
			repo := sessionrepo.New(test.conn)
			ctx := context.Background()
			s := entities.Session{
				UserId:       u.ID,
				RefreshToken: "123",
				UserAgent:    "123",
				Ip:           "10.01.01.01",
				ExpiresIn:    time.Now().AddDate(0, 0, 30),
			}

			if test.doError {
				conn.Close()
			}

			_, err := repo.New(ctx, s)
			assert.Equal(t, test.expectedError, err != nil)
		})
	}
}

func TestRefreshRepo_Get(t *testing.T) {
	conn, cleaner := testhelpers.RealDb()
	defer cleaner("users", "refresh_sessions")

	repo := user2.NewUserRepo(conn)
	user := entities.User{Company: "123"}
	_ = user.GenerateSecrets()

	u, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	sesRepo := sessionrepo.New(conn)
	session := entities.Session{
		UserId:       u.ID,
		RefreshToken: "123",
		UserAgent:    "123",
		Ip:           "123",
		ExpiresIn:    time.Now().AddDate(0, 0, 30),
	}
	s, err := sesRepo.New(context.Background(), session)
	assert.NoError(t, err)

	var tt = []struct {
		name          string
		conn          connector.Conn
		token         string
		expectedError bool
		doError       bool
	}{
		{
			name:          "mock | should get sessionrepo from db",
			conn:          mockRefreshGetFuncConnSuccess{},
			token:         "123",
			expectedError: false,
			doError:       false,
		},
		{
			name:          "should get sessionrepo from db with token = " + s.RefreshToken,
			conn:          conn,
			token:         "123",
			expectedError: false,
			doError:       false,
		},
		{
			name:          "mock | get error if sessionrepo not found",
			conn:          mockRefreshGetFuncConnError{},
			token:         "net tot token =)",
			expectedError: true,
			doError:       false,
		},
		{
			name:          "get error if sessionrepo not found",
			conn:          conn,
			token:         "net tot token =)",
			expectedError: true,
			doError:       false,
		},
		{
			name:          "get error if conn closed",
			conn:          conn,
			token:         "nu a tyt i ne nado =)",
			expectedError: true,
			doError:       true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			repo := sessionrepo.New(test.conn)
			ctx := context.Background()

			if test.doError {
				conn.Close()
			}

			_, err := repo.Get(ctx, test.token)
			assert.Equal(t, test.expectedError, err != nil)
		})
	}
}

func TestRefreshRepo_Remove(t *testing.T) {
	conn, cleaner := testhelpers.RealDb()
	defer cleaner("users", "refresh_sessions")

	repo := user2.NewUserRepo(conn)
	user := entities.User{Company: "123"}
	_ = user.GenerateSecrets()

	u, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	sesRepo := sessionrepo.New(conn)
	session := entities.Session{
		UserId:       u.ID,
		RefreshToken: "123",
		UserAgent:    "123",
		Ip:           "123",
		ExpiresIn:    time.Now().AddDate(0, 0, 30),
	}
	_, err = sesRepo.New(context.Background(), session)
	assert.NoError(t, err)

	var tt = []struct {
		name          string
		conn          connector.Conn
		token         string
		expectedError bool
		doError       bool
	}{
		{
			name:          "mock | should remove sessionrepo from db",
			conn:          mockRefreshGetFuncConnSuccess{},
			token:         "123",
			expectedError: false,
			doError:       false,
		},
		{
			name:          "mock | get error if sessionrepo not found",
			conn:          mockRefreshGetFuncConnError{},
			token:         "123",
			expectedError: true,
			doError:       false,
		},
		{
			name:          "should remove sessionrepo with token = 123",
			conn:          conn,
			token:         "123",
			expectedError: false,
			doError:       false,
		},
		{
			name:          "get error if sessionrepo not found",
			conn:          conn,
			token:         "123",
			expectedError: true,
			doError:       false,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			repo := sessionrepo.New(test.conn)
			ctx := context.Background()

			if test.doError {
				conn.Close()
			}

			_, err := repo.Remove(ctx, test.token)
			assert.Equal(t, test.expectedError, err != nil)
		})
	}
}

func TestRefreshRepo_RemoveBatch(t *testing.T) {
	conn, cleaner := testhelpers.RealDb()
	defer cleaner("users", "refresh_sessions")

	repo := user2.NewUserRepo(conn)
	user := entities.User{Company: "123"}
	_ = user.GenerateSecrets()

	u, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	sesRepo := sessionrepo.New(conn)
	session := entities.Session{
		UserId:       u.ID,
		RefreshToken: "123",
		UserAgent:    "123",
		Ip:           "123",
		ExpiresIn:    time.Now().AddDate(0, 0, 30),
	}
	ids := make([]int, 3)
	for i := 0; i < 3; i++ {
		s, err := sesRepo.New(context.Background(), session)
		assert.NoError(t, err)
		session.RefreshToken += fmt.Sprint(i)
		ids[i] = s.ID
	}

	var tt = []struct {
		name          string
		conn          connector.Conn
		ids           []int
		expectedError bool
	}{
		{
			name:          "mock | should remove sessionrepo from db",
			conn:          mockRefreshBatchFuncConnSuccess{},
			ids:           []int{1, 2, 3},
			expectedError: false,
		},
		{
			name:          "should remove sessionrepo with ids",
			conn:          conn,
			ids:           ids,
			expectedError: false,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			repo := sessionrepo.New(test.conn)
			ctx := context.Background()

			err := repo.RemoveBatch(ctx, test.ids...)
			assert.Equal(t, test.expectedError, err != nil)
		})
	}
}