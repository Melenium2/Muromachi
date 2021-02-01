package store

import (
	"Muromachi/store/connector"
	"Muromachi/store/sessions"
	"Muromachi/store/userrepo"
)

type AuthCollection struct {
	Sessions sessions.Sessions
	Users    userrepo.UsersRepo
}

func NewAuthCollection(conn connector.Conn) *AuthCollection {
	return &AuthCollection{
		Sessions: nil,
		Users:    userrepo.NewUserRepo(conn),
	}
}
