package userstore_test

import (
	"Muromachi/config"
	"Muromachi/store/entities"
	"Muromachi/store/testhelpers"
	user2 "Muromachi/store/users/userstore"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepo_Create_ShouldCreateNewUserAndPutItToDatabase_Mock(t *testing.T) {
	conn := mockUserConnectionSuccess{}
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	var (
		err  error
		user = entities.User{
			Company: "123",
		}
	)
	_ = user.GenerateSecrets()

	user, err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.Greater(t, user.ID, 0)
	assert.NotEmpty(t, user.ClientId)
	assert.NotEmpty(t, user.ClientSecret)
}

func TestUserRepo_Create_ShouldCreateNewUserAndPutItToDatabase(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users")
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	var (
		err  error
		user = entities.User{
			Company: "123",
		}
	)
	_ = user.GenerateSecrets()

	user, err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.Greater(t, user.ID, 0)
	assert.NotEmpty(t, user.ClientId)
	assert.NotEmpty(t, user.ClientSecret)
}

func TestUserRepo_Create_ShouldReturnErrorBecauseUserHasNotClientIdAndSecret_Mock(t *testing.T) {
	conn := mockUserConnectionSuccess{}
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	var (
		user = entities.User{
			Company: "123",
		}
	)

	_, err := repo.Create(ctx, user)
	assert.Error(t, err)
}

func TestUserRepo_Create_ShouldReturnErrorBecauseUserHasNotClientIdAndSecret(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users")
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	var (
		err  error
		user = entities.User{
			Company: "123",
		}
	)

	_, err = repo.Create(ctx, user)
	assert.Error(t, err)
}

func TestUserRepo_Create_ShouldReturnErrorIfCanNotCreateUser_Mock(t *testing.T) {
	conn := mockUserConnectionError{}
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	var (
		user = entities.User{
			Company: "123",
		}
	)

	_, err := repo.Create(ctx, user)
	assert.Error(t, err)
}

func TestUserRepo_Create_ShouldReturnErrorIfCanNotCreateUser(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users")
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	conn.Close()
	var (
		err  error
		user = entities.User{
			Company: "123",
		}
	)
	_ = user.GenerateSecrets()

	_, err = repo.Create(ctx, user)
	assert.Error(t, err)
}

func TestUserRepo_Approve_ShouldGetUserFromDatabase_Mock(t *testing.T) {
	conn := mockUserApproveConnectionSuccess{}
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	usr, err := repo.Approve(ctx, "nudopustim1")
	assert.NoError(t, err)
	assert.Greater(t, usr.ID, 0)
	assert.NotEmpty(t, usr.ClientSecret)
	assert.NotEmpty(t, usr.ClientId)
}

func TestUserRepo_Approve_ShouldGetUserFromDatabase(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, cleaner := testhelpers.RealDb(cfg.Database)
	defer cleaner("users")
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	var (
		err  error
		user = entities.User{Company: "123"}
	)
	_ = user.GenerateSecrets()

	secret := user.ClientSecret

	user, err = repo.Create(ctx, user)
	assert.NoError(t, err)

	usr, err := repo.Approve(ctx, user.ClientId)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, usr.ID)
	assert.Equal(t, user.ClientId, usr.ClientId)

	assert.NoError(t, usr.CompareSecret(secret))
}

func TestUserRepo_Approve_ShouldReturnErrorIfClientIdDoesNotExist_Mock(t *testing.T) {
	conn := mockUserApproveConnectionError{}
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	_, err := repo.Approve(ctx, "nudopustim1")
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
}

func TestUserRepo_Approve_ShouldReturnErrorIfClientIdDoesNotExist(t *testing.T) {
	cfg := config.New("../../../config/dev.yml")
	cfg.Database.Schema = "../../../config/schema.sql"

	conn, _ := testhelpers.RealDb(cfg.Database)
	repo := user2.NewUserRepo(conn)
	ctx := context.Background()

	_, err := repo.Approve(ctx, "nudopustim1")
	assert.Error(t, err)
	assert.Equal(t, pgx.ErrNoRows, err)
}
