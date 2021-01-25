package config_test

import (
	"Muromachi/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_ShouldCreateValidInstanceOfDatabaseConfig_Dev(t *testing.T) {
	cfg := config.New()

	assert.NotEmpty(t, cfg.Database.Database)
	assert.NotEmpty(t, cfg.Database.User)
	assert.NotEmpty(t, cfg.Database.Password)
	assert.NotEmpty(t, cfg.Database.Schema)
	assert.NotEmpty(t, cfg.Database.Address)
	assert.NotEmpty(t, cfg.Database.Port)
}

func TestConfig_ShouldCreateValidInstanceOfDatabaseConfig_Prod(t *testing.T) {
	cfg := config.New("./prod.yml")

	assert.NotEmpty(t, cfg.Database.Database)
	assert.NotEmpty(t, cfg.Database.Schema)
	assert.NotEmpty(t, cfg.Database.Port)

	assert.Empty(t, cfg.Database.User)
	assert.Empty(t, cfg.Database.Password)
	assert.Empty(t, cfg.Database.Address)
}