package store

import (
	"context"
	"time"
)

type UserRepo struct {
	conn Conn
}

// Creating new user. For creation only need company name
//
// The returned user contains the client id and secret for api requests
func (u *UserRepo) Create(ctx context.Context, user User) (User, error) {
	var (
		oldSecret string
		err       error
	)
	// Using bcrypt for hashing user secret in database
	if oldSecret, err = user.SecureSecret(); err != nil {
		return User{}, err
	}
	if user.AddedAt.IsZero() {
		user.AddedAt = time.Now().UTC()
	}
	row := u.conn.QueryRow(
		ctx,
		"insert into users (clientId, clientSecret, company, addedAt) values ($1, $2, $3, $4) returning id",
		user.ClientId, user.ClientSecret, user.Company, user.AddedAt,
	)
	var id int
	if err = row.Scan(&id); err != nil {
		return User{}, err
	}

	user.ID = id
	user.ClientSecret = oldSecret
	return user, nil
}

// Search for a user by clientId and get it
func (u *UserRepo) Approve(ctx context.Context, clientId string) (User, error) {
	row := u.conn.QueryRow(
		ctx,
		"select * from users where clientId = $1",
		clientId,
	)
	var user User
	if err := row.Scan(
		&user.ID,
		&user.ClientId,
		&user.ClientSecret,
		&user.Company,
		&user.AddedAt); err != nil {
		return User{}, err
	}

	return user, nil
}

func NewUserRepo(conn Conn) *UserRepo {
	return &UserRepo{
		conn: conn,
	}
}
