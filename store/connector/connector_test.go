package connector_test

import (
	"Muromachi/config"
	"Muromachi/store/connector"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var c = config.New("../../config/dev.yml")

func TestConnectionUrl_ShouldCreateValidUrlFromConfig(t *testing.T) {
	url, err := connector.ConnectionUrl(c.Database)

	assert.NoError(t, err)
	assert.Equal(t, "postgresql://postgres:123456@localhost:5432/tracking", url)
}

func TestConnectionUrl_ShouldReturnDefaultDatabaseIfNameIsNotProvided(t *testing.T) {
	cfg := config.New("../../config/dev.yml")
	cfg.Database.Database = ""
	url, err := connector.ConnectionUrl(cfg.Database)

	assert.NoError(t, err)
	assert.Equal(t, "postgresql://postgres:123456@localhost:5432/default", url)
}

func TestConnectionUrl_ShouldReturnErrorIfAddressNotProvided(t *testing.T) {
	cfg := config.New("../../config/dev.yml")
	cfg.Database.Address = ""
	_, err := connector.ConnectionUrl(cfg.Database)
	assert.Error(t, err)
}

func TestConnectionUrl_ShouldReturnErrorIfPortNotProvided(t *testing.T) {
	cfg := config.New("../../config/dev.yml")
	cfg.Database.Port = ""
	_, err := connector.ConnectionUrl(cfg.Database)
	assert.Error(t, err)
}

func TestConnect_ShouldEstablishInitialConnection(t *testing.T) {
	url, err := connector.ConnectionUrl(c.Database)
	assert.NoError(t, err)

	log.Print(url)
	conn, err := connector.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestConnect_ShouldReturnErrorIfCanNotEstablishInitialConnection(t *testing.T) {
	cfg := config.New("../../config/dev.yml")
	url, _ := connector.ConnectionUrl(cfg.Database)
	url += "123???123123123"
	conn, err := connector.Connect(url)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestInitSchema_ShouldCreateDatabaseSchemaOrDoNothing(t *testing.T) {
	url, _ := connector.ConnectionUrl(c.Database)

	conn, err := connector.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	if err == nil {
		c.Database.Schema = "../../config/schema.sql"
		err = connector.InitSchema(conn, c.Database.Schema)
		assert.NoError(t, err)
	}
}

func TestInitSchema_ShouldReturnErrorIfFilepathWrong(t *testing.T) {
	cfg := config.New("../../config/dev.yml")
	url, _ := connector.ConnectionUrl(cfg.Database)

	conn, err := connector.Connect(url)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	cfg.Database.Schema = "../../config/schema"
	err = connector.InitSchema(conn, cfg.Database.Schema)
	assert.Error(t, err)
}
