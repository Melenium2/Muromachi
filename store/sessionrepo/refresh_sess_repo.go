package sessionrepo

import (
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type RefreshSessions interface {
	New(ctx context.Context, session entities.Session) (entities.Session, error)
	Get(ctx context.Context, token string) (entities.Session, error)
	Remove(ctx context.Context, token string) (entities.Session, error)
	RemoveBatch(ctx context.Context, sessionid ...int) error
	UserSessions(ctx context.Context, userId int) ([]entities.Session, error)
}

type RefreshRepo struct {
	conn connector.Conn
}

func (r *RefreshRepo) New(ctx context.Context, session entities.Session) (entities.Session, error) {
	row := r.conn.QueryRow(
		ctx,
		"insert into refresh_sessions (userId, refreshToken, useragent, ip, expiresIn) values ($1, $2, $3, $4, $5) returning id, createdAt",
		session.UserId, session.RefreshToken, session.UserAgent, session.Ip, session.ExpiresIn,
	)
	var id int
	var t time.Time
	if err := row.Scan(&id, &t); err != nil {
		return entities.Session{}, err
	}

	session.ID = id
	session.CreatedAt = t

	return session, nil
}

func (r *RefreshRepo) Get(ctx context.Context, token string) (entities.Session, error) {
	row := r.conn.QueryRow(
		ctx,
		"select * from refresh_sessions where refreshToken = $1",
		token,
	)
	var session entities.Session
	if err := row.Scan(
		&session.ID,
		&session.UserId,
		&session.RefreshToken,
		&session.UserAgent,
		&session.Ip,
		&session.ExpiresIn,
		&session.CreatedAt,
	); err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (r *RefreshRepo) Remove(ctx context.Context, token string) (entities.Session, error) {
	row := r.conn.QueryRow(
		ctx,
		"delete from refresh_sessions where refreshToken = $1 returning *",
		token,
	)
	var session entities.Session
	if err := row.Scan(
		&session.ID,
		&session.UserId,
		&session.RefreshToken,
		&session.UserAgent,
		&session.Ip,
		&session.ExpiresIn,
		&session.CreatedAt,
	); err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (r *RefreshRepo) RemoveBatch(ctx context.Context, sessionid ...int) error {
	if len(sessionid) == 0 {
		return nil
	}
	var ids string
	for _, v := range sessionid {
		ids += fmt.Sprintf("%d,", v)
	}
	ids = ids[:len(ids)-1]
	_, err := r.conn.Exec(
		ctx,
		"delete from refresh_sessions where id in (" + ids + ")",
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshRepo) UserSessions(ctx context.Context, userId int) ([]entities.Session, error) {
	var sess entities.Session
	var sessions []entities.Session

	_, err := r.conn.QueryFunc(
		ctx,
		"select * from refresh_sessions where userId = $1",
		[]interface{} { userId },
		[]interface{} {
			&sess.ID, &sess.UserId, &sess.RefreshToken, &sess.UserAgent,
			&sess.Ip, &sess.ExpiresIn, &sess.CreatedAt,
		},
		func(row pgx.QueryFuncRow) error {
			sessions = append(sessions, sess)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func New(conn connector.Conn) *RefreshRepo {
	return &RefreshRepo{
		conn: conn,
	}
}
