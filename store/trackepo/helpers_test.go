package trackepo_test

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"time"
)

// Mock connection with errors (Keyword and Category tables)
type mockTrackConnectionErrors struct {
}

func (m mockTrackConnectionErrors) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockTrackConnectionErrors) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockTrackErrorRow{}
}

func (m mockTrackConnectionErrors) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, pgx.ErrNoRows
}

type mockTrackErrorRow struct {
}

func (mr mockTrackErrorRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockTrackErrorRow) RawValues() [][]byte {
	return nil
}

func (mr mockTrackErrorRow) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}

// Mock connection with successful returned objects (Keyword or Category tables)
type mockTrackConnection struct {
}

func (m mockTrackConnection) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (m mockTrackConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	t, _ := time.Parse("2006-01-02", "2021-01-19")
	return mockTrackRow{
		Id:             78,
		BundleId:       12,
		Type:           "type",
		Date:           t,
		Place:          18,
		AppId:          12,
		AppBundle:      "123",
		AppCategory:    "FINANCE",
		AppDeveloperId: "com.develoeper",
		AppDeveloper:   "super developer",
		AppGeo:         "ru_RU",
		AppStartAt:     t.AddDate(-1, 0, 0),
		AppPeriod:      31,
	}
}

func (m mockTrackConnection) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	t, _ := time.Parse("2006-01-02", "2021-01-19")
	for i := 0; i < 4; i++ {
		_ = mockTrackRow{
			Id:             78,
			BundleId:       12,
			Type:           "type",
			Date:           t,
			Place:          18,
			AppId:          12,
			AppBundle:      "123",
			AppCategory:    "FINANCE",
			AppDeveloperId: "com.develoeper",
			AppDeveloper:   "super developer",
			AppGeo:         "ru_RU",
			AppStartAt:     t.AddDate(-1, 0, 0),
			AppPeriod:      31,
		}.Scan(scans...)
		t = t.AddDate(0, 0, 1)
		_ = f(mockTrackRow{})
	}

	return nil, nil
}

type mockTrackRow struct {
	Id             int       `json:"-"`
	BundleId       int       `json:"bundle,omitempty"`
	Type           string    `json:"type,omitempty"`
	Date           time.Time `json:"date,omitempty"`
	Place          int32     `json:"place,omitempty"`
	AppId          int       `json:"-"`
	AppBundle      string    `json:"appbundle,omitempty"`
	AppCategory    string    `json:"category,omitempty"`
	AppDeveloperId string    `json:"developer_id,omitempty"`
	AppDeveloper   string    `json:"developer,omitempty"`
	AppGeo         string    `json:"geo,omitempty"`
	AppStartAt     time.Time `json:"start_at,omitempty"`
	AppPeriod      uint32    `json:"period,omitempty"`
}

func (mr mockTrackRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockTrackRow) RawValues() [][]byte {
	return nil
}

func (mr mockTrackRow) Scan(dest ...interface{}) error {
	Id := dest[0].(*int)
	BundleId := dest[1].(*int)
	Type := dest[2].(*string)
	Place := dest[3].(*int32)
	Date := dest[4].(*time.Time)
	AppId := dest[5].(*int)
	AppBundle := dest[6].(*string)
	AppCategory := dest[7].(*string)
	AppDeveloperId := dest[8].(*string)
	AppDeveloper := dest[9].(*string)
	AppGeo := dest[10].(*string)
	AppStartAt := dest[11].(*time.Time)
	AppPeriod := dest[12].(*uint32)

	*Id = mr.Id
	*BundleId = mr.BundleId
	*Type = mr.Type
	*Date = mr.Date
	*Place = mr.Place
	*AppId = mr.AppId
	*AppBundle = mr.AppBundle
	*AppCategory = mr.AppCategory
	*AppDeveloperId = mr.AppDeveloperId
	*AppDeveloper = mr.AppDeveloper
	*AppGeo = mr.AppGeo
	*AppStartAt = mr.AppStartAt
	*AppPeriod = mr.AppPeriod

	return nil
}
