package server_test

import (
	"Muromachi/store/entities"
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

type mockSession struct {
}

func (m mockSession) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m mockSession) CheckIfExist(ctx context.Context, key string) error {
	return nil
}

func (m mockSession) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	session.ID = 1
	return session, nil
}

func (m mockSession) Get(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSession) Remove(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{
		ID: 1,
	}, nil
}

func (m mockSession) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return nil
}

func (m mockSession) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	return nil, nil
}

type mockConn struct {

}

func (m mockConn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	id := args[0].(string)
	if id != "123" {
		return mockRowError{}
	}

	user := entities.User{
		ID:           123,
		ClientId:     "123",
		ClientSecret: "123",
		Company:      "123",
		AddedAt:      time.Now(),
	}
	_, _ = user.SecureSecret()
	return mockRow{
		ID:           user.ID,
		ClientId:     user.ClientId,
		ClientSecret: user.ClientSecret,
		Company:      user.Company,
		AddedAt:      user.AddedAt,
	}
}

func (m mockConn) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockConn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

type mockRow struct {
	ID           int       `json:"id,omitempty"`
	ClientId     string    `json:"client_id,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	Company      string    `json:"company,omitempty"`
	AddedAt      time.Time `json:"added_at,omitempty"`
}

func (m mockRow) Scan(dest ...interface{}) error {
	ID := dest[0].(*int)
	ClientId := dest[1].(*string)
	ClientSecret := dest[2].(*string)
	Company := dest[3].(*string)
	AddedAt := dest[4].(*time.Time)

	*ID = m.ID
	*ClientId = m.ClientId
	*ClientSecret = m.ClientSecret
	*Company = m.Company
	*AddedAt = m.AddedAt

	return nil
}

type mockRowError struct {
}

func (m mockRowError) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}


