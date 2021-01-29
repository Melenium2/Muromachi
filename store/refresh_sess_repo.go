package store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type RefreshRepo struct {
	conn Conn
}

func (r *RefreshRepo) New(ctx context.Context, session Session) (Session, error) {
	row := r.conn.QueryRow(
		ctx,
		"insert into refresh_sessions (userId, refreshToken, useragent, ip, expiresIn) values ($1, $2, $3, $4, $5) returning id, createdAt",
		session.UserId, session.RefreshToken, session.UserAgent, session.Ip, session.ExpiresIn,
	)
	var id int
	var t time.Time
	if err := row.Scan(&id, &t); err != nil {
		return Session{}, err
	}

	session.ID = id
	session.CreatedAt = t

	return session, nil
}

func (r *RefreshRepo) Get(ctx context.Context, token string) (Session, error) {
	row := r.conn.QueryRow(
		ctx,
		"select * from refresh_sessions where refreshToken = $1",
		token,
	)
	var session Session
	if err := row.Scan(
		&session.ID,
		&session.UserId,
		&session.RefreshToken,
		&session.UserAgent,
		&session.Ip,
		&session.ExpiresIn,
		&session.CreatedAt,
	); err != nil {
		return Session{}, err
	}

	return session, nil
}

func (r *RefreshRepo) Remove(ctx context.Context, token string) (Session, error) {
	row := r.conn.QueryRow(
		ctx,
		"delete from refresh_sessions where refreshToken = $1 returning *",
		token,
	)
	var session Session
	if err := row.Scan(
		&session.ID,
		&session.UserId,
		&session.RefreshToken,
		&session.UserAgent,
		&session.Ip,
		&session.ExpiresIn,
		&session.CreatedAt,
	); err != nil {
		return Session{}, err
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
	row := r.conn.QueryRow(
		ctx,
		"delete from refresh_sessions where id in (" + ids + ")",
	)
	if err := row.Scan();err != nil {
		return err
	}

	return nil
}

func (r *RefreshRepo) UserSessions(ctx context.Context, userId int) ([]Session, error) {
	var sess Session
	var sessions []Session

	_, err := r.conn.QueryFunc(
		ctx,
		"select * from refresh_sessions where userId = $1",
		[]interface{} { userId },
		[]interface{} {
			&sess.UserId, &sess.RefreshToken, &sess.UserAgent,
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

func NewRefreshRepo(conn Conn) *RefreshRepo {
	return &RefreshRepo{
		conn: conn,
	}
}
