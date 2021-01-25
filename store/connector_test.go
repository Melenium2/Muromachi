package store_test

import (
	"Muromachi/config"
	"Muromachi/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

var c = config.New("../config/dev.yml")

func TestConnectionUrl_ShouldCreateValidUrlFromConfig(t *testing.T) {
	url, err := store.ConnectionUrl(c.Database)

	assert.NoError(t, err)
	assert.Equal(t, "postgresql://postgres:123456@localhost:5432/tracking", url)
}

func TestConnect_ShouldEstablishInitialConnection(t *testing.T) {
	url, _ := store.ConnectionUrl(c.Database)

	conn, err := store.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestInitSchema_ShouldCreateDatabaseSchemaOrDoNothing(t *testing.T) {
	url, _ := store.ConnectionUrl(c.Database)

	conn, err := store.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	c.Database.Schema = "../config/schema.sql"
	err = store.InitSchema(conn, c.Database.Schema)
	assert.NoError(t, err)
}
