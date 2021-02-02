package store

import (
	"Muromachi/store/apprepo"
	"Muromachi/store/connector"
	"Muromachi/store/metarepo"
	"Muromachi/store/trackepo"
)

type TableCollection struct {
	App  TrackingRepo
	Meta TrackingRepo
	Cat  TrackingRepo
	Keys TrackingRepo
}

func NewTrackingCollection(conn connector.Conn) *TableCollection {
	return &TableCollection{
		App:  apprepo.New(conn),
		Meta: metarepo.New(conn),
		Cat:  trackepo.NewCat(conn),
		Keys: trackepo.NewKeys(conn),
	}
}
