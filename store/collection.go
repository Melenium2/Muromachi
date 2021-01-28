package store

type TableCollection struct {
	App  TrackingRepo
	Meta TrackingRepo
	Cat  TrackingRepo
	Keys TrackingRepo
}

func New(conn Conn) *TableCollection {
	return &TableCollection{
		App:  NewApp(conn),
		Meta: NewMeta(conn),
		Cat:  NewCat(conn),
		Keys: NewKeys(conn),
	}
}

type AuthCollection struct {
	Sessions Sessions
	Users    UsersRepo
}
