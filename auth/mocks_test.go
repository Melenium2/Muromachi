package auth_test

import (
	"Muromachi/store/entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type mockSession struct {
}

func (m mockSession) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m mockSession) CheckIfExist(ctx context.Context, key string) error {
	return fmt.Errorf("%s", "key not found")
}

func (m mockSession) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	session.ID = 1
	return session, nil
}

func (m mockSession) Get(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{
		ID:           1,
		UserId:       123,
		RefreshToken: "123",
		UserAgent:    "123",
		Ip:           "10.10.0.1",
		ExpiresIn:    time.Now().AddDate(0, 0, 1),
		CreatedAt:    time.Now(),
	}, nil
}

func (m mockSession) Remove(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{
		ID:           1,
		UserId:       123,
		RefreshToken: "123",
		UserAgent:    "123",
		Ip:           "10.10.0.1",
		ExpiresIn:    time.Now().AddDate(0, 0, 1),
		CreatedAt:    time.Now(),
	}, nil
}

func (m mockSession) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return nil
}

func (m mockSession) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	return []entities.Session{
		{
			ID:           1,
			UserId:       123,
			RefreshToken: "123",
			UserAgent:    "123",
			Ip:           "10.10.0.1",
			ExpiresIn:    time.Now().AddDate(0, 0, 1),
			CreatedAt:    time.Now(),
		},
	}, nil
}

type mockSessionRemoveNoRows struct {
}

func (m mockSessionRemoveNoRows) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m mockSessionRemoveNoRows) CheckIfExist(ctx context.Context, key string) error {
	return nil
}

func (m mockSessionRemoveNoRows) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionRemoveNoRows) Get(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionRemoveNoRows) Remove(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{}, pgx.ErrNoRows
}

func (m mockSessionRemoveNoRows) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return nil
}

func (m mockSessionRemoveNoRows) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	return nil, nil
}

type mockSessionRemoveExpiredSession struct {
}

func (m mockSessionRemoveExpiredSession) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m mockSessionRemoveExpiredSession) CheckIfExist(ctx context.Context, key string) error {
	return nil
}

func (m mockSessionRemoveExpiredSession) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionRemoveExpiredSession) Get(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionRemoveExpiredSession) Remove(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{
		ExpiresIn: time.Now().Add(time.Hour * -2),
	}, nil
}

func (m mockSessionRemoveExpiredSession) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return nil
}

func (m mockSessionRemoveExpiredSession) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	return nil, nil
}

type mockSessionMoreThen5 struct {
}

func (m mockSessionMoreThen5) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m mockSessionMoreThen5) CheckIfExist(ctx context.Context, key string) error {
	return nil
}

func (m mockSessionMoreThen5) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	session.ID = 1
	return session, nil
}

func (m mockSessionMoreThen5) Get(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionMoreThen5) Remove(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionMoreThen5) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return nil
}

func (m mockSessionMoreThen5) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	sessions := make([]entities.Session, 0)
	for i := 0; i < 7; i++ {
		sessions = append(sessions, entities.Session{ID: i+1})
	}
	return sessions, nil
}

type mockSessionBannedToken struct {

}

func (m mockSessionBannedToken) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m mockSessionBannedToken) CheckIfExist(ctx context.Context, key string) error {
	// This means that the session is not in the black list
	return nil
}

func (m mockSessionBannedToken) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionBannedToken) Get(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{}, nil
}

func (m mockSessionBannedToken) Remove(ctx context.Context, token string) (entities.Session, error) {
	return entities.Session{
		ID:           1,
		UserId:       123,
		RefreshToken: "123",
		UserAgent:    "123",
		Ip:           "10.10.0.1",
		ExpiresIn:    time.Now().AddDate(0, 0, 1),
		CreatedAt:    time.Now(),
	}, nil
}

func (m mockSessionBannedToken) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return nil
}

func (m mockSessionBannedToken) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	return nil, nil
}

