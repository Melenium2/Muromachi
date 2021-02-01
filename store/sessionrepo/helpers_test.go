package sessionrepo_test

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

//
// Method: New
// Case: success
//
type mockRefreshNewFuncRowSuccess struct {
	Id        int
	CreatedAt time.Time
}

func (m mockRefreshNewFuncRowSuccess) Scan(dest ...interface{}) error {
	id := dest[0].(*int)
	createdAt := dest[1].(*time.Time)

	*id = m.Id
	*createdAt = m.CreatedAt

	return nil
}

type mockRefreshNewFuncConnSuccess struct {
}

func (m mockRefreshNewFuncConnSuccess) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockRefreshNewFuncConnSuccess) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRefreshNewFuncRowSuccess{
		Id:        1,
		CreatedAt: time.Now().UTC(),
	}
}

func (m mockRefreshNewFuncConnSuccess) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

//
// Method: New
// Case: Error
//

type mockRefreshNewFuncRowError struct {
}

func (m mockRefreshNewFuncRowError) Scan(dest ...interface{}) error {
	return fmt.Errorf("%s", "can not create sessionrepo")
}

type mockRefreshNewFuncConnError struct {
}

func (m mockRefreshNewFuncConnError) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockRefreshNewFuncConnError) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRefreshNewFuncRowError{}
}

func (m mockRefreshNewFuncConnError) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

//
// Method: Get, Remove
// Case: Success
//

type mockRefreshGetFuncRowSuccss struct {
	ID           int       `json:"id,omitempty"`
	UserId       int       `json:"user_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Ip           string    `json:"ip,omitempty"`
	ExpiresIn    time.Time `json:"expires_in,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

func (m mockRefreshGetFuncRowSuccss) Scan(dest ...interface{}) error {
	id := dest[0].(*int)
	userid := dest[1].(*int)
	refreshToken:= dest[2].(*string)
	ua := dest[3].(*string)
	ip := dest[4].(*string)
	expires := dest[5].(*time.Time)
	created := dest[6].(*time.Time)

	*id = m.ID
	*userid = m.UserId
	*refreshToken = m.RefreshToken
	*ua = m.UserAgent
	*ip = m.Ip
	*expires = m.ExpiresIn
	*created = m.CreatedAt

	return nil
}

type mockRefreshGetFuncConnSuccess struct {

}

func (m mockRefreshGetFuncConnSuccess) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockRefreshGetFuncConnSuccess) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRefreshGetFuncRowSuccss{
		ID:           283,
		UserId:       16,
		RefreshToken: "123",
		UserAgent:    "123",
		Ip:           "123",
		ExpiresIn:    time.Now(),
		CreatedAt:    time.Now(),
	}
}

func (m mockRefreshGetFuncConnSuccess) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

//
// Method: Get, Remove
// Case: Error
//

type mockRefreshGetFuncRowError struct {
}

func (m mockRefreshGetFuncRowError) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}

type mockRefreshGetFuncConnError struct {

}

func (m mockRefreshGetFuncConnError) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockRefreshGetFuncConnError) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRefreshGetFuncRowError{}
}

func (m mockRefreshGetFuncConnError) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

//
// Method: RemoveBatch
// Case: Success
//

type mockRefreshBatchFuncRowSuccess struct {

}

func (m mockRefreshBatchFuncRowSuccess) Scan(dest ...interface{}) error {
	return nil
}

type mockRefreshBatchFuncConnSuccess struct {

}

func (m mockRefreshBatchFuncConnSuccess) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockRefreshBatchFuncConnSuccess) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRefreshBatchFuncRowSuccess{}
}

func (m mockRefreshBatchFuncConnSuccess) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

//
// Method: RemoveBatch
// Case: Error
//

type mockRefreshBatchFuncRowError struct {

}

func (m mockRefreshBatchFuncRowError) Scan(dest ...interface{}) error {
	return nil
}

type mockRefreshBatchFuncConnError struct {

}

func (m mockRefreshBatchFuncConnError) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockRefreshBatchFuncConnError) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRefreshBatchFuncRowError{}
}

func (m mockRefreshBatchFuncConnError) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, nil
}

