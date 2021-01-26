package store_test

import (
	"Muromachi/config"
	"Muromachi/store"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var c = config.New("../config/dev.yml")

func TestConnectionUrl_ShouldCreateValidUrlFromConfig(t *testing.T) {
	url, err := store.ConnectionUrl(c.Database)

	assert.NoError(t, err)
	assert.Equal(t, "postgresql://postgres:123456@localhost:5432/tracking", url)
}

func TestConnectionUrl_ShouldReturnDefaultDatabaseIfNameIsNotProvided(t *testing.T) {
	cfg := config.New("../config/dev.yml")
	cfg.Database.Database = ""
	url, err := store.ConnectionUrl(cfg.Database)

	assert.NoError(t, err)
	assert.Equal(t, "postgresql://postgres:123456@localhost:5432/default", url)
}

func TestConnectionUrl_ShouldReturnErrorIfAddressNotProvided(t *testing.T) {
	cfg := config.New("../config/dev.yml")
	cfg.Database.Address = ""
	_, err := store.ConnectionUrl(cfg.Database)
	assert.Error(t, err)
}

func TestConnectionUrl_ShouldReturnErrorIfPortNotProvided(t *testing.T) {
	cfg := config.New("../config/dev.yml")
	cfg.Database.Port = ""
	_, err := store.ConnectionUrl(cfg.Database)
	assert.Error(t, err)
}

func TestConnect_ShouldEstablishInitialConnection(t *testing.T) {
	url, err := store.ConnectionUrl(c.Database)
	assert.NoError(t, err)

	log.Print(url)
	conn, err := store.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestConnect_ShouldReturnErrorIfCanNotEstablishInitialConnection(t *testing.T) {
	cfg := config.New("../config/dev.yml")
	url, _ := store.ConnectionUrl(cfg.Database)
	url += "123???123123123"
	conn, err := store.Connect(url)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestInitSchema_ShouldCreateDatabaseSchemaOrDoNothing(t *testing.T) {
	url, _ := store.ConnectionUrl(c.Database)

	conn, err := store.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	if err == nil {
		c.Database.Schema = "../config/schema.sql"
		err = store.InitSchema(conn, c.Database.Schema)
		assert.NoError(t, err)
	}
}

func TestInitSchema_ShouldReturnErrorIfFilepathWrong(t *testing.T) {
	cfg := config.New("../config/dev.yml")
	url, _ := store.ConnectionUrl(cfg.Database)

	conn, err := store.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	cfg.Database.Schema = "../config/schema"
	err = store.InitSchema(conn, cfg.Database.Schema)
	assert.Error(t, err)
}
