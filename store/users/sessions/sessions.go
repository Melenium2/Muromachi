package sessions

import (
	"Muromachi/store/entities"
	"Muromachi/store/users/sessions/blacklist"
	"Muromachi/store/users/sessions/tokens"
	"context"
	"time"
)

// Main interface to working with user sessions
type Session interface {
	blacklist.BlackList
	tokens.RefreshSession
}

type sessionsImpl struct {
	sessions  tokens.RefreshSession
	blacklist blacklist.BlackList
}

// Add session to black list
func (s sessionsImpl) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return s.blacklist.Add(ctx, key, value, ttl)
}

// Check if sessions already exists in black list
func (s sessionsImpl) CheckIfExist(ctx context.Context, key string) error {
	return s.blacklist.CheckIfExist(ctx, key)
}

// Remove session from black list
func (s sessionsImpl) Del(ctx context.Context, keys ...string) (int64, error) {
	return s.blacklist.Del(ctx, keys...)
}

// Create new session
func (s sessionsImpl) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	return s.sessions.New(ctx, session)
}

// Get existing session
func (s sessionsImpl) Get(ctx context.Context, token string) (entities.Session, error) {
	return s.sessions.Get(ctx, token)
}

// Remove session
func (s sessionsImpl) Remove(ctx context.Context, token string) (entities.Session, error) {
	return s.sessions.Remove(ctx, token)
}

// Remove batch sessions
func (s sessionsImpl) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return s.sessions.RemoveBatch(ctx, sessionid...)
}

// Get user sessions
func (s sessionsImpl) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	return s.sessions.UserSessions(ctx, userId)
}

func New(sessions tokens.RefreshSession, blacklist blacklist.BlackList) *sessionsImpl {
	return &sessionsImpl{
		sessions:  sessions,
		blacklist: blacklist,
	}
}
