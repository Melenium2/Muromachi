package entities_test

import (
	"Muromachi/store/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_GenerateSecrets_ShouldGenerateNewSecretsWithoutErrors(t *testing.T) {
	user := entities.User{
		Company:      "Random name",
		AddedAt:      time.Now(),
	}

	assert.NoError(t, user.GenerateSecrets())
	assert.NotEmpty(t, user.ClientSecret)
	assert.NotEmpty(t, user.ClientId)
	t.Log(user.ClientId)
	t.Log(user.ClientSecret)
}

func TestUser_SecureSecret_ShouldGenerateSecureHashToClientSecret(t *testing.T) {
	user := entities.User{
		Company:      "Random name",
		AddedAt:      time.Now(),
	}
	assert.NoError(t, user.GenerateSecrets())
	secret := user.ClientSecret

	_, err := user.SecureSecret()
	assert.NoError(t, err)
	assert.NotEqual(t, secret, user.ClientSecret)
}

func TestUser_CompareSecret_ShouldValidateUserByComparingSecrets(t *testing.T) {
	user := entities.User{
		Company:      "Random name",
		AddedAt:      time.Now(),
	}
	assert.NoError(t, user.GenerateSecrets())
	secret := user.ClientSecret
	_, err := user.SecureSecret()
	assert.NoError(t, err)
	assert.NoError(t, user.CompareSecret(secret))
}

