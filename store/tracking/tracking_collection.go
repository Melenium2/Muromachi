package tracking

import (
	"Muromachi/store/connector"
)

// helper to working with different tables
type Tables struct {
	App  Repository
	Meta Repository
	Cat  Repository
	Keys Repository
}

func NewTrackingTables(conn connector.Conn) *Tables {
	return &Tables{
		App:  NewAppRepo(conn),
		Meta: NewMetaRepo(conn),
		Cat:  NewCatRepo(conn),
		Keys: NewKeysRepo(conn),
	}
}
