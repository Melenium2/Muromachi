package blacklist_test

import (
	"Muromachi/config"
	"Muromachi/store/testhelpers"
	"Muromachi/store/users/sessions/blacklist"
	"Muromachi/utils"
	"context"
	"fmt"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlackList_Add_Mock(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := blacklist.New(db)
	ctx := context.Background()

	hash := utils.Hash("123", "123")
	mock.ExpectSetNX(hash, 123, time.Minute*2).SetVal(true)
	assert.NoError(t, repo.Add(ctx, hash, 123, time.Minute*2))
}

func TestBlackList_Add(t *testing.T) {
	cfg := config.New("../../../../config/dev.yml")

	conn, cleaner := testhelpers.RedisDb(cfg.Database.Redis)
	defer cleaner()
	repo := blacklist.New(conn)
	ctx := context.Background()

	var tt = []struct {
		name          string
		ttl           time.Duration
		doError       bool
		expectedError bool
	}{
		{
			name:          "should return right hash from redis",
			ttl:           time.Minute * 2,
			doError:       false,
			expectedError: false,
		},
		{
			name:          "should return error if hash not found",
			ttl:           time.Minute * 2,
			doError:       true,
			expectedError: true,
		},
		{
			name:          "should return error if ttl is expired",
			ttl:           time.Second * 1,
			doError:       false,
			expectedError: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			hash := utils.Hash("123", "123")
			assert.NoError(t, repo.Add(ctx, hash, 123, test.ttl))

			if test.doError {
				hash += "???123"
			}
			// test ttl
			time.Sleep(time.Second * 3)

			_, err := conn.Get(ctx, hash).Int()
			assert.Equal(t, err != nil, test.expectedError)
		})
	}
}

func TestBlackList_CheckIfExist_Mock_ShouldGetValue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := blacklist.New(db)
	ctx := context.Background()

	hash := utils.Hash("123", "123")
	mock.ExpectGet(hash).SetVal(hash)
	err := repo.CheckIfExist(ctx, hash)
	assert.NoError(t, err)
}

func TestBlackList_CheckIfExist_Mock_ShouldReturnNil(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := blacklist.New(db)
	ctx := context.Background()

	hash := utils.Hash("123", "123")
	mock.ExpectGet(hash).RedisNil()
	err := repo.CheckIfExist(ctx, hash)
	assert.Error(t, err)
}

func TestBlackList_CheckIfExist_Mock_ShouldReturnUnexpectedError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := blacklist.New(db)
	ctx := context.Background()

	hash := utils.Hash("123", "123")
	mock.ExpectGet(hash).SetErr(fmt.Errorf("%s", "unexpected error"))
	err := repo.CheckIfExist(ctx, hash)
	assert.Error(t, err)
}

func TestBlackList_CheckIfExist(t *testing.T) {
	cfg := config.New("../../../../config/dev.yml")

	conn, cleaner := testhelpers.RedisDb(cfg.Database.Redis)
	defer cleaner()
	repo := blacklist.New(conn)
	ctx := context.Background()

	hash := utils.Hash("123", "123")
	_ = repo.Add(ctx, hash, 123, time.Minute*3)

	var tt = []struct {
		name          string
		hash          string
		expectedError bool
	}{
		{
			name:          "should get value without errors",
			hash:          hash,
			expectedError: false,
		},
		{
			name:          "should receive error if key not found",
			hash:          "123",
			expectedError: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			err := repo.CheckIfExist(ctx, test.hash)
			assert.Equal(t, err != nil, test.expectedError)
		})
	}
}
