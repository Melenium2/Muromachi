package sessions

import (
	"Muromachi/store/banrepo"
	"Muromachi/store/entities"
	"Muromachi/store/refreshrepo"
	"context"
)

type Sessions interface {
	banrepo.BlackList
	refreshrepo.RefreshSessions
}

type sessionsImpl struct {
	sessions  refreshrepo.RefreshSessions
	blacklist banrepo.BlackList
}

func (s sessionsImpl) AddBlock() {
}

func (s sessionsImpl) CheckBlock() {
}

func (s sessionsImpl) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	return s.sessions.New(ctx, session)
}

func (s sessionsImpl) Get(ctx context.Context, token string) (entities.Session, error) {
	return s.sessions.Get(ctx, token)
}

func (s sessionsImpl) Remove(ctx context.Context, token string) (entities.Session, error) {
	return s.sessions.Remove(ctx, token)
}

func (s sessionsImpl) RemoveBatch(ctx context.Context, sessionid ...int) error {
	return s.sessions.RemoveBatch(ctx, sessionid...)
}

func (s sessionsImpl) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	return s.sessions.UserSessions(ctx, userId)
}

func New(sessions refreshrepo.RefreshSessions, blacklist banrepo.BlackList) *sessionsImpl {
	return &sessionsImpl{
		sessions: sessions,
	}
}
