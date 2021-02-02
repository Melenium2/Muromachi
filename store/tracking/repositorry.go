package tracking

import (
	"Muromachi/store/connector"
	"Muromachi/store/entities"
	"Muromachi/store/tracking/appstore"
	"Muromachi/store/tracking/metastore"
	"Muromachi/store/tracking/trackstore"
	"context"
	"time"
)

type Repository interface {
	ProducerFunc(ctx context.Context, sql string, params ...interface{}) (entities.DboSlice, error)
	ByBundleId(ctx context.Context, bundleId int) (entities.DboSlice, error)
	TimeRange(ctx context.Context, bundleId int, start, end time.Time) (entities.DboSlice, error)
	LastUpdates(ctx context.Context, bundleId, count int) (entities.DboSlice, error)
}

func NewCatRepo(conn connector.Conn) *trackstore.CatRepo {
	return &trackstore.CatRepo{
		Conn: conn,
	}
}

func NewKeysRepo(conn connector.Conn) *trackstore.KeysRepo {
	return &trackstore.KeysRepo{
		Conn: conn,
	}
}


func NewMetaRepo(conn connector.Conn) *metastore.Repo {
	return &metastore.Repo{
		Conn: conn,
	}
}

func NewAppRepo(conn connector.Conn) *appstore.Repo {
	return &appstore.Repo{
		Conn: conn,
	}
}