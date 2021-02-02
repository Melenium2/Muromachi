package connector

import (
	"Muromachi/config"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"strings"
)

// ConnectionUrl creates database connection url by given config.DBConfig
func ConnectionUrl(config config.DBConfig) (string, error) {
	url := "postgresql://"

	if config.User != "" && config.Password != "" {
		url += config.User + ":" + config.Password + "@"
	}

	if config.Address == "" {
		return "", errors.New("empty db address")
	}
	url += config.Address

	if config.Port == "" {
		return "", errors.New("empty db port")
	}
	url += ":" + config.Port + "/"

	db := config.Database
	if db == "" {
		db = "default"
	}
	url += db

	return url, nil
}

// Connect connect to database by url
func Connect(url string) (*pgxpool.Pool, error) {
	connect, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return connect, nil
}

// InitSchema creates database schema from shemafile string
func InitSchema(connection *pgxpool.Pool, schemafile string) error {
	b, err := ioutil.ReadFile(schemafile)
	if err != nil {
		return err
	}
	schema := string(b)
	tables := strings.Split(schema, ";\r\n")

	ctx := context.Background()
	for _, v := range tables {
		_, err = connection.Exec(ctx, v)
		if err != nil {
			if strings.Contains(err.Error(), "developercontacts") && strings.Contains(err.Error(), "42710")  {
				continue
			}
			return err
		}
	}


	return nil
}

func EstablishPostgresConnection(config config.DBConfig) (*pgxpool.Pool, error) {
	url, err := ConnectionUrl(config)
	if err != nil {
		return nil, err
	}

	conn, err := Connect(url)
	if err != nil {
		return nil, err
	}

	err = InitSchema(conn, config.Schema)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func EstablishRedisConnection(config config.RedisConfig) (*redis.Client, error) {
	// TODO Сделать коннектор для редиса
	return nil, nil
}
