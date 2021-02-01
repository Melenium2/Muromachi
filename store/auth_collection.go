package store

import (
	"Muromachi/store/connector"
	"Muromachi/store/userrepo"
)

type AuthCollection struct {
	Sessions Sessions
	Users    userrepo.UsersRepo
}

func NewAuthCollection(conn connector.Conn) *AuthCollection {
	return &AuthCollection{
		Sessions: nil,
		Users:    userrepo.NewUserRepo(conn),
	}
}
