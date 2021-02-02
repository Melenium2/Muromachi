package userstore_test

import (
	"Muromachi/store/entities"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

//
// Create method
//
type mockUserRowSuccess struct {
	ID int `json:"id,omitempty"`
}

func (m mockUserRowSuccess) Scan(dest ...interface{}) error {
	id := dest[0].(*int)
	*id = m.ID

	return nil
}

type mockUserConnectionSuccess struct {
}

func (m mockUserConnectionSuccess) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockUserConnectionSuccess) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockUserRowSuccess{
		ID: 1,
	}
}

func (m mockUserConnectionSuccess) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

type mockUserRowError struct {
}

func (m mockUserRowError) Scan(dest ...interface{}) error {
	return fmt.Errorf("%s", "can not save userrepo")
}

type mockUserConnectionError struct {
}

func (m mockUserConnectionError) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockUserConnectionError) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockUserRowError{}
}

func (m mockUserConnectionError) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

//
// Approve method
//
type mockUserApproveRowSuccess struct {
	ID           int       `json:"id,omitempty"`
	ClientId     string    `json:"client_id,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	Company      string    `json:"company,omitempty"`
	AddedAt      time.Time `json:"added_at,omitempty"`
}

func (m mockUserApproveRowSuccess) Scan(dest ...interface{}) error {
	id := dest[0].(*int)
	clientID := dest[1].(*string)
	clientSecret := dest[2].(*string)
	company := dest[3].(*string)
	addedAt := dest[4].(*time.Time)

	*id = m.ID
	*clientID = m.ClientId
	*clientSecret = m.ClientSecret
	*company = m.Company
	*addedAt = m.AddedAt

	return nil
}

type mockUserApproveConnectionSuccess struct {
}

func (m mockUserApproveConnectionSuccess) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockUserApproveConnectionSuccess) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	user := entities.User{Company: "123"}
	_ = user.GenerateSecrets()
	_, _ = user.SecureSecret()
	return mockUserApproveRowSuccess{
		ID:           1,
		ClientId:     user.ClientId,
		ClientSecret: user.ClientSecret,
		Company:      user.Company,
		AddedAt:      time.Now(),
	}
}

func (m mockUserApproveConnectionSuccess) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

type mockUserApproveRowError struct {
}

func (m mockUserApproveRowError) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}

type mockUserApproveConnectionError struct {
}

func (m mockUserApproveConnectionError) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockUserApproveConnectionError) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockUserApproveRowError{}
}

func (m mockUserApproveConnectionError) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

