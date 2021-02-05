package users

import (
	"Muromachi/store/users/sessions"
	"Muromachi/store/users/userstore"
)

// helper to working with different tables
type Tables struct {
	Sessions sessions.Session
	Users    userstore.UsersRepo
}

func NewAuthTables(session sessions.Session, userRepo userstore.UsersRepo) *Tables {
	return &Tables{
		Sessions: session,
		Users:    userRepo,
	}
}
