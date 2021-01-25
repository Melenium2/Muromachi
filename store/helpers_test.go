package store_test

import (
	"Muromachi/config"
	"Muromachi/store"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"log"
	"strings"
	"time"
)

func RealDb() (*pgx.Conn, func(names ...string)) {
	c := config.New("../../../config/dev.yml")
	url, err := store.ConnectionUrl(c.Database)
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		panic(err)
	}

	if err := store.InitSchema(conn, "../../../config/schema.sql"); err != nil {
		panic(err)
	}

	return conn, func(names ...string) {
		_, err := conn.Exec(context.Background(), fmt.Sprintf("truncate table %s CASCADE", strings.Join(names, ",")))
		if err != nil {
			log.Print(err)
		}
	}
}

// Mock connection with errors
type mockAppConnectionErrors struct {
}

func (m mockAppConnectionErrors) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return mockRowError{}
}

func (m mockAppConnectionErrors) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	return nil, pgx.ErrNoRows
}

type mockRowError struct {
}

func (mre mockRowError) Scan(dest ...interface{}) error {
	return pgx.ErrNoRows
}

func (mre mockRowError) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mre mockRowError) RawValues() [][]byte {
	return nil
}

// Mock connection with successful returned objects
type mockAppConnection struct {
}

func (m mockAppConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	t, _ := time.Parse("2006-01-01", "2020-01-01")
	return mockRow{
		Id:          1,
		Bundle:      "com.test",
		Category:    "FINANCE",
		DeveloperId: "92834848476158744",
		Developer:   "Random",
		Geo:         "ru_ru",
		StartAt:     t,
		Period:      30,
	}
}

func (m mockAppConnection) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	t, _ := time.Parse("2006-01-01", "2020-01-01")
	for i := 0; i < 2; i++ {
		err := mockRow{
			Id:          i+1,
			Bundle:      "com.test",
			Category:    "FINANCE",
			DeveloperId: "92834848476158744",
			Developer:   "Random",
			Geo:         "ru_ru",
			StartAt:     t,
			Period:      30,
		}.Scan(scans...)
		if err != nil {
			return nil, err
		}
		err = f(mockRow{})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type mockRow struct {
	Id          int
	Bundle      string
	Category    string
	DeveloperId string
	Developer   string
	Geo         string
	StartAt     time.Time
	Period      uint32
}

func (mr mockRow) FieldDescriptions() []pgproto3.FieldDescription {
	return nil
}

func (mr mockRow) RawValues() [][]byte {
	return nil
}

func (mr mockRow) Scan(dest ...interface{}) error {
	Id := dest[0].(*int)
	Bundle := dest[1].(*string)
	Category := dest[2].(*string)
	DeveloperId := dest[3].(*string)
	Developer := dest[4].(*string)
	Geo := dest[5].(*string)
	StartAt := dest[6].(*time.Time)
	Period := dest[7].(*uint32)

	*Id = mr.Id
	*Bundle = mr.Bundle
	*Category = mr.Category
	*DeveloperId = mr.DeveloperId
	*Developer = mr.Developer
	*Geo = mr.Geo
	*StartAt = mr.StartAt
	*Period = mr.Period

	return nil
}
